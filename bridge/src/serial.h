struct Response {
  char returnKey[16];
  char returnValue[32];
  bool hasReturn;
  bool ok;
  int errCode;
  char errName[32];
  char errDetails[128];
};

Response execute(const char*, const char* = 0);
