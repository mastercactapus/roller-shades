#include <ESP8266SSDP.h>

#include "settings.h"

void initSSDP() {
  SSDP.setHTTPPort(80);
  SSDP.setSchemaURL("description.xml");
  SSDP.setDeviceType("urn:diy-org:device:DIYRollerShade:1");
  SSDP.setName(settings.name);
  SSDP.begin();
}