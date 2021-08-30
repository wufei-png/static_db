package localcache

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.sz.sensetime.com/viper/infra/infra-model-manager/api"
)

var (
	ErrFileLocked   = errors.New("model file locked")
	ErrFileNotFound = errors.New("model file not found")

	ErrNoMatchingBlob = errors.New("no matching blob")
)

const (
	blobDir = "blobs"
	metaDir = "meta"
)

type LocalCache struct {
	rootDir string
}

func NewLocalCache(rootDir string) (*LocalCache, error) {
	if rootDir == "" {
		return nil, errors.New("invalid root dir")
	}
	log.Debug("initializing local cache: ", rootDir)
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(rootDir, blobDir), 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(rootDir, metaDir), 0755); err != nil {
		return nil, err
	}
	return &LocalCache{
		rootDir: rootDir,
	}, nil
}

func (c *LocalCache) removeBadMetaFile(fp string, err error) {
	log.Warn("invalid metadata: ", err, ", removing ", fp)
	if err := os.Remove(fp); err != nil {
		log.Error("failed to remove bad metafile: ", err)
	}
}

func (c *LocalCache) Get(mpath *api.ModelPath) (*MetaFile, string, error) {
	fp, err := ModelPathToFilePath(mpath, false)
	if err != nil {
		return nil, "", nil
	}
	fp = filepath.Join(c.rootDir, metaDir, fp)

	model, err := NewMetaFileFromPath(fp)
	if err != nil {
		if err != ErrFileNotFound {
			c.removeBadMetaFile(fp, err)
		}
		return nil, "", err
	}
	blobPath, err := validateBlob(c.rootDir, model.Meta)
	if err != nil {
		return nil, "", err
	}

	return model, blobPath, nil
}

func (c *LocalCache) BlobReader(oid, checksum string) (io.ReadCloser, error) {
	blobPath, err := GetBlobPath(oid, checksum)
	if err != nil {
		return nil, err
	}
	blobPath = filepath.Join(c.rootDir, blobDir, blobPath)
	return os.Open(blobPath)
}

func (c *LocalCache) ValidateBlob(meta *api.Model) (string, error) {
	return validateBlob(c.rootDir, meta)
}

func (c *LocalCache) LockForWriter(meta *MetaFile) (*FileWriter, error) {
	return newFileWriter(c.rootDir, meta)
}

func (c *LocalCache) CopyLocalFile(model *api.Model, path string) error {
	checksum, size, err := ChecksumFile(path)
	if err != nil {
		return err
	}
	model.Checksum = checksum
	model.Size = size
	meta, err := NewMetaFileFromModel(model)
	if err != nil {
		return err
	}
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close() // nolint

	w, err := c.LockForWriter(meta)
	if err != nil {
		return err
	}
	err = CopyWithChecksum(w, r, size, checksum)
	if err != nil {
		w.Abort() // nolint
		return err
	}
	return w.Commit()
}

func validateBlob(rootDir string, model *api.Model) (string, error) {
	blobPath, err := GetBlobPath(model.GetOid(), model.GetChecksum())
	if err != nil {
		return "", err
	}
	blobPath = filepath.Join(rootDir, blobDir, blobPath)
	fi, err := os.Stat(blobPath)
	if err != nil {
		// log.Warn("no matching blob for oid: ", blobPath)
		return "", ErrNoMatchingBlob
	}
	if fi.IsDir() || fi.Size() != model.GetSize() {
		log.Warn("no matching blob size for oid: ", blobPath)
		return "", ErrNoMatchingBlob
	}

	return blobPath, nil
}

func (c *LocalCache) RemoveMeta(model *api.Model) error {
	metaPath, err := ModelPathToFilePath(model.GetModelPath(), false)
	if err != nil {
		return err
	}
	metaPath = filepath.Join(c.rootDir, metaDir, metaPath)
	return os.Remove(metaPath)
}

func (c *LocalCache) ListModels() ([]*MetaFile, error) {
	var models []*MetaFile
	base := filepath.Join(c.rootDir, metaDir)
	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warn("error during walk: ", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		model, err := NewMetaFileFromPath(path)
		if err != nil {
			log.Warn("invalid meta content: ", err)
			return nil
		}

		models = append(models, model)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (c *LocalCache) CommitMetaFile(meta *MetaFile) error {
	log.Info("committing model meta only: ", meta.HashKey, ", ", meta.Meta)
	fp, err := ModelPathToFilePath(meta.Meta.GetModelPath(), false)
	if err != nil {
		return nil
	}
	fp = filepath.Join(c.rootDir, metaDir, fp)
	if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return err
	}

	tmp := fp + TempSuffix()
	metaWriter, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	if _, err := metaWriter.Write(meta.Data); err != nil {
		metaWriter.Close() // nolint
		return err
	}
	if err := metaWriter.Sync(); err != nil {
		metaWriter.Close() // nolint
		return err
	}
	if err := metaWriter.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp, fp); err != nil {
		return err
	}

	return nil
}

// CleanupBlob removes blob files which make isOrphan return true
// and remove empty dirs
func (c *LocalCache) CleanupBlobs(isOrphan func(blobPath string, modifiedAt time.Time) bool) error {
	blobPath := filepath.Join(c.rootDir, blobDir)
	err := filepath.Walk(blobPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warn("error during walk: ", err)
			return err
		}
		if info.IsDir() {
			return nil
		}

		if isOrphan(filepath.Base(path), info.ModTime()) {
			if err = os.Remove(path); err != nil {
				log.Warn("failed to remove ", path, err)
				return nil
			}
			log.Info("unreferenced blob: ", path, " removed, modified at: ", info.ModTime())
		}
		return nil
	})

	return err
}

// StorageStatus returns disk utility
func (c *LocalCache) StorageStatus() (*api.StorageStatus, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(c.rootDir, &fs)
	if err != nil {
		return nil, err
	}

	total := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize)
	available := fs.Bavail * uint64(fs.Bsize)
	used := total - free
	return &api.StorageStatus{
		Total:     total,
		Used:      used,
		Available: available,
	}, nil
}
