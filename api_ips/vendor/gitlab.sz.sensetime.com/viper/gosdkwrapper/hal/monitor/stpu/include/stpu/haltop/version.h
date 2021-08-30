#ifndef __HALTOP_VERSION_H__
#define __HALTOP_VERSION_H__

#ifdef __cplusplus
extern "C" const char *haltopVersion(void);
extern "C" const char *haltopRevision(void);
#else
extern const char *haltopVersion(void);
extern const char *haltopRevision(void);
#endif

#endif