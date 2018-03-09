#include <ESP8266mDNS.h>

#include "mdns.h"
#include "settings.h"

void initMDNS() {
  if (!MDNS.begin(sanitizeDNS(settings.name))) {
    Serial.println("Error setting up MDNS responder!");
  }
  MDNS.addService("http", "tcp", 80);
  MDNS.addService("diy-roller-shade", "tcp", 80);
  MDNS.setInstanceName(settings.name);

  const char* cname = (const char*)settings.name;
  MDNS.addServiceTxt("diy-roller-shade", "tcp", "name", cname);
}

const char* sanitizeDNS(const char* name) {
  static char safeName[64];
  int p = 0;
  for (int i = 0; i < 256; i++) {
    char c = name[i];
    if (c == 0) {
      safeName[p] = 0;
      break;
    }
    if (p == 63) {
      safeName[p] = 0;
      break;
    }

    if (c == '-') {
      // can't start with hyphen
      if (p == 0) {
        p++;
        continue;
      }
      // can't have double hyphen
      if (safeName[p - 1] == '-') {
        p++;
        continue;
      }
    } else if (c >= 'A' && c <= 'Z') {
      // make lower case
      c += 32;
    } else if ((c < '0') || c > 'z' || (c > '9' && c < 'a')) {
      // skip non alphanumeric
      if (p > 0 && safeName[p - 1] != '-') {
        safeName[p] = '-';
      }
      p++;
      continue;
    }
    safeName[p] = c;
    p++;
  }
  if (p == 0) {
    strcpy(safeName, "diy-roller-shades");
    return safeName;
  }
  if (safeName[p - 1] == '-') {
    safeName[p - 1] = '-';
  }
  return safeName;
}
