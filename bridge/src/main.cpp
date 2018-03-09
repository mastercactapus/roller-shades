
#include <ESP8266WiFi.h>
#include <elapsedMillis.h>

#include "events.h"
#include "http.h"
#include "mdns.h"
#include "serial.h"
#include "settings.h"
#include "ssdp.h"
#include "wifi-creds.h"

#define SERIAL_RATE 115200
#define SERIAL_TIMEOUT 1000

void setup() {
  Serial.begin(SERIAL_RATE);
  Serial.setTimeout(SERIAL_TIMEOUT);
  Serial.println();
  WiFi.begin(ssid, pass);
  saveSettings();

  Serial.print("Connecting");
  while (WiFi.status() != WL_CONNECTED) {
    Serial.print(".");
    delay(500);
  }
  Serial.println();

  Serial.print("Connected, IP address: ");
  Serial.println(WiFi.localIP());

  initHTTPServer();
  initMDNS();
  initSSDP();
}
void serialFlush() {
  while (Serial.available() > 0) {
    char t = Serial.read();
  }
}

void loop() {
  serialFlush();
  HTTP.handleClient();

  static char oldEventState[256] = "";
  static elapsedMillis sinceLastPush = 0;
  // send events up to once per sec
  if (sinceLastPush > 1000) {
    updateEventState();
    if (strlen(settings.eventURL) && strcmp(eventState, oldEventState) != 0) {
      sendEvents(settings.eventURL);
      strcpy(oldEventState, eventState);
    }
    sinceLastPush = 0;
  }
}
