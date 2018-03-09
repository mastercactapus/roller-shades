#include <Arduino.h>

#include "serial.h"

Response execute(const char* command, const char* arg) {
  static Response res;
  Serial.println();
  if (arg) {
    Serial.print(command);
    Serial.println(arg);
  } else {
    Serial.println(command);
  }

  String resp = Serial.readStringUntil('\n');
  resp.trim();
  int delim;
  if (resp.startsWith("+")) {
    res.hasReturn = true;
    delim = resp.indexOf(':');
    resp.substring(1, delim).toCharArray(res.returnKey, sizeof(res.returnKey));
    resp.substring(delim + 1).toCharArray(res.returnValue,
                                          sizeof(res.returnValue));
    resp = Serial.readStringUntil('\n');
    resp.trim();
  } else {
    res.hasReturn = false;
    res.returnValue[0] = 0;
    res.returnKey[0] = 0;
  }

  res.ok = resp.equals("OK");
  if (res.ok) {
    res.errCode = 200;
    res.errDetails[0] = 0;
    res.errName[0] = 0;
    return res;
  }
  if (!resp.startsWith("ERR:")) {
    res.errCode = 500;
    strcpy(res.errName, "Serial Failure");
    strcpy(res.errDetails, "Serial communication failure.");
    return res;
  }

  delim = resp.indexOf(':');
  resp = resp.substring(delim + 1);
  delim = resp.indexOf(':');
  res.errCode = resp.substring(0, delim).toInt();
  resp = resp.substring(delim + 1);
  delim = resp.indexOf(':');
  resp.substring(0, delim).toCharArray(res.errName, sizeof(res.errName));
  resp = resp.substring(delim + 1);
  resp.toCharArray(res.errDetails, sizeof(res.errDetails));

  return res;
}
