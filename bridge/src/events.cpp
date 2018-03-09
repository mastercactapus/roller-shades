#include <ESP8266HTTPClient.h>

#include "serial.h"
#include "settings.h"

char eventState[256] = "";

void updateEventState() {
  char pos[4] = "";
  char lock[2] = "";
  char status[16] = "";

  strcpy(pos, execute("AT+POS?").returnValue);
  strcpy(pos, execute("AT+LOCK?").returnValue);
  strcpy(pos, execute("AT+STATUS?").returnValue);

  sprintf(eventState, "pos=%s&lock=%s&status=%s", pos, lock, status);
}

void sendEvents(const char *url) {
  HTTPClient http;
  http.setReuse(false);
  http.setTimeout(1000);
  http.begin(url);
  http.addHeader("Content-Type", "application/x-www-form-urlencoded");
  http.POST(eventState);
  http.end();
}
