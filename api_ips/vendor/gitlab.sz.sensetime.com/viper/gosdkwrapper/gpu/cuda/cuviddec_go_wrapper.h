CUresult wrapCuvidGetDecoderCaps(CUVIDDECODECAPS *pdc) {
    return cuvidGetDecoderCaps(pdc);
}
  
CUresult wrapCuvidCreateDecoder(CUvideodecoder *phDecoder, CUVIDDECODECREATEINFO *pdci) {
    return cuvidCreateDecoder(phDecoder ,pdci);
}
  
CUresult wrapCuvidDestroyDecoder(CUvideodecoder hDecoder) {
    return cuvidDestroyDecoder(hDecoder);
}
  
CUresult wrapCuvidDecodePicture(CUvideodecoder hDecoder, CUVIDPICPARAMS *pPicParams) {
    return cuvidDecodePicture(hDecoder ,pPicParams);
}
  
CUresult wrapCuvidUnmapVideoFrame(CUvideodecoder hDecoder, unsigned int DevPtr) {
    return cuvidUnmapVideoFrame(hDecoder ,DevPtr);
}
  
CUresult wrapCuvidUnmapVideoFrame64(CUvideodecoder hDecoder, unsigned long long DevPtr) {
    return cuvidUnmapVideoFrame64(hDecoder ,DevPtr);
}
  
CUresult wrapCuvidCtxLockCreate(CUvideoctxlock *pLock, CUcontext ctx) {
    return cuvidCtxLockCreate(pLock ,ctx);
}
  
CUresult wrapCuvidCtxLockDestroy(CUvideoctxlock lck) {
    return cuvidCtxLockDestroy(lck);
}
  
CUresult wrapCuvidCtxLock(CUvideoctxlock lck, unsigned int reserved_flags) {
    return cuvidCtxLock(lck ,reserved_flags);
}
  
CUresult wrapCuvidCtxUnlock(CUvideoctxlock lck, unsigned int reserved_flags) {
    return cuvidCtxUnlock(lck ,reserved_flags);
}
  
