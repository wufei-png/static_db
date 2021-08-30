#ifndef __HAL_ML_VERSION_H__
#define __HAL_ML_VERSION_H__

#ifdef __cplusplus
extern "C" const char *stpumlVersion(void);
extern "C" const char *stpumlRevision(void);
#else
extern const char *stpumlVersion(void);
extern const char *stpumlRevision(void);
#endif

#endif