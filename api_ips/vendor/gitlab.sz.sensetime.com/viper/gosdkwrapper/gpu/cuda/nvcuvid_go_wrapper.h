CUresult wrapCuvidCreateVideoSource(CUvideosource *pObj, const char *pszFileName, CUVIDSOURCEPARAMS *pParams) {
    return cuvidCreateVideoSource(pObj ,pszFileName ,pParams);
}
  
CUresult wrapCuvidCreateVideoSourceW(CUvideosource *pObj, const wchar_t *pwszFileName, CUVIDSOURCEPARAMS *pParams) {
    return cuvidCreateVideoSourceW(pObj ,pwszFileName ,pParams);
}
  
CUresult wrapCuvidDestroyVideoSource(CUvideosource obj) {
    return cuvidDestroyVideoSource(obj);
}
  
CUresult wrapCuvidSetVideoSourceState(CUvideosource obj, cudaVideoState state) {
    return cuvidSetVideoSourceState(obj ,state);
}
  
CUresult wrapCuvidGetVideoSourceState(CUvideosource obj) {
    return cuvidGetVideoSourceState(obj);
}
  
CUresult wrapCuvidGetSourceVideoFormat(CUvideosource obj, CUVIDEOFORMAT *pvidfmt, unsigned int flags) {
    return cuvidGetSourceVideoFormat(obj ,pvidfmt ,flags);
}
  
CUresult wrapCuvidGetSourceAudioFormat(CUvideosource obj, CUAUDIOFORMAT *paudfmt, unsigned int flags) {
    return cuvidGetSourceAudioFormat(obj ,paudfmt ,flags);
}
  
CUresult wrapCuvidCreateVideoParser(CUvideoparser *pObj, CUVIDPARSERPARAMS *pParams) {
    return cuvidCreateVideoParser(pObj ,pParams);
}
  
CUresult wrapCuvidParseVideoData(CUvideoparser obj, CUVIDSOURCEDATAPACKET *pPacket) {
    return cuvidParseVideoData(obj ,pPacket);
}
  
CUresult wrapCuvidDestroyVideoParser(CUvideoparser obj) {
    return cuvidDestroyVideoParser(obj);
}
  
