package localcache

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type FileWriter struct {
	meta *MetaFile

	flock *fileLock

	blobPath string
	metaPath string

	blobWriter *os.File
	rootDir    string

	tmpSuffix string
}

func newFileWriter(rootDir string, meta *MetaFile) (*FileWriter, error) {
	fp, err := ModelPathToFilePath(meta.Meta.GetModelPath(), false)
	if err != nil {
		return nil, err
	}
	fp = filepath.Join(rootDir, metaDir, fp)
	if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return nil, err
	}

	blobPath, err := GetBlobPath(meta.Meta.GetOid(), meta.Meta.GetChecksum())
	if err != nil {
		return nil, err
	}
	blobPath = filepath.Join(rootDir, blobDir, blobPath)
	if err := os.MkdirAll(filepath.Dir(blobPath), 0755); err != nil {
		return nil, err
	}

	flockPath, err := ModelPathToFilePath(meta.Meta.GetModelPath(), true)
	if err != nil {
		return nil, err
	}
	flockPath = filepath.Join(rootDir, metaDir, flockPath)

	flock := newFileLock(flockPath)
	if err := flock.TryLock(); err != nil {
		return nil, err
	}

	tmpSuffix := TempSuffix()

	log.Debug("file lock acquired: ", flockPath)
	blobWriter, err := os.OpenFile(blobPath+tmpSuffix, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		if err := flock.Unlock(); err != nil {
			log.Warn("failed to unlock: ", err)
		}
		return nil, err
	}

	return &FileWriter{
		meta:       meta,
		flock:      flock,
		blobPath:   blobPath,
		metaPath:   fp,
		blobWriter: blobWriter,
		rootDir:    rootDir,
		tmpSuffix:  tmpSuffix,
	}, nil
}

func (w *FileWriter) Write(b []byte) (int, error) {
	return w.blobWriter.Write(b)
}

func (w *FileWriter) unlock() {
	if err := w.flock.Unlock(); err != nil {
		log.Warn("failed to unlock flock: ", err)
	}
	if err := w.flock.Remove(); err != nil {
		log.Warn("failed to remove flock: ", err)
	}
	log.Debug("file lock released: ", w.flock.GetFilename())
}

func (w *FileWriter) removeTmp() {
	if err := os.Remove(w.blobPath + w.tmpSuffix); err != nil {
		log.Warn("failed to remove temp: ", err)
	}
}

func (w *FileWriter) Commit() error {
	defer w.unlock()

	if err := w.blobWriter.Sync(); err != nil {
		log.Error("failed to sync file: ", err)
		w.blobWriter.Close() // nolint
		w.removeTmp()
		return err
	}
	if err := w.blobWriter.Close(); err != nil {
		log.Error("failed to close file: ", err)
		w.removeTmp()
		return err
	}

	fi, err := os.Stat(w.blobPath + w.tmpSuffix)
	if err != nil {
		w.removeTmp()
		return err
	}
	if fi.Size() != w.meta.Meta.GetSize() {
		log.Error("blob file ", w.blobPath, " size mismatch: ", fi.Size(), " vs. ", w.meta.Meta.GetSize())
		w.removeTmp()
		return ErrNoMatchingBlob
	}

	metaWriter, err := os.OpenFile(w.metaPath+w.tmpSuffix, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	if _, err := metaWriter.Write(w.meta.Data); err != nil {
		metaWriter.Close() // nolint
		w.removeTmp()
		return err
	}
	if err := metaWriter.Sync(); err != nil {
		metaWriter.Close() // nolint
		w.removeTmp()
		return err
	}
	if err := metaWriter.Close(); err != nil {
		w.removeTmp()
		return err
	}

	log.Info("committing model: ", w.metaPath, ", meta: ", w.meta.Meta)
	if err := os.Rename(w.blobPath+w.tmpSuffix, w.blobPath); err != nil {
		w.removeTmp()
		return err
	}
	if err := os.Rename(w.metaPath+w.tmpSuffix, w.metaPath); err != nil {
		// os.Remove(w.blobPath)
		return err
	}
	return nil
}

func (w *FileWriter) Abort() error {
	defer w.unlock()
	log.Info("abort model writer: ", w.metaPath)
	if err := w.blobWriter.Close(); err != nil {
		log.Warn("failed to close: ", err)
	}
	if err := os.Remove(w.blobPath + w.tmpSuffix); err != nil {
		log.Warn("failed to remove: ", err)
	}
	return nil
}
