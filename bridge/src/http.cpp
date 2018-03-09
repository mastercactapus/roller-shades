#include <ESP8266SSDP.h>
#include <ESP8266WebServer.h>
#include <FS.h>

#include "http.h"
#include "serial.h"
#include "settings.h"

ESP8266WebServer HTTP(80);

void initHTTPServer() {
  SPIFFS.begin();
  HTTP.onNotFound([]() {
    if (!handleFileRead(HTTP.uri())) {
      HTTP.send(404, "text/plain", "Not found.");
    }
  });

  HTTP.on("/api/events", HTTP_GET,
          []() { HTTP.send(200, "text/plain", settings.eventURL); });
  HTTP.on("/api/events", HTTP_DELETE, []() {
    settings.eventURL[0] = 0;
    saveSettings();
    HTTP.send(204, "text/plain", "");
  });
  HTTP.on("/api/events", HTTP_POST, []() {
    HTTP.arg("url").toCharArray(settings.eventURL, 256);
    saveSettings();
    HTTP.send(200, "text/plain", "");
  });

  HTTP.on("/description.xml", HTTP_GET, []() { SSDP.schema(HTTP.client()); });
  HTTP.on("/api/name", HTTP_GET,
          []() { HTTP.send(200, "text/plain", settings.name); });
  HTTP.on("/api/name", HTTP_POST, []() {
    String newName = HTTP.arg("name");
    newName.trim();
    if (newName.length() > 128) {
      HTTP.send(400, "text/plain",
                "Name too long (must be less than 128 bytes).");
      return;
    }
    for (byte i = 0; i < newName.length(); i++) {
      byte c = newName[i];
      if (c < ' ' || c > '~') {
        HTTP.send(400, "text/plain", "Name contains invalid characters.");
        return;
      }
    }

    strcpy(settings.name, newName.c_str());
    saveSettings();

    HTTP.send(200, "text/plain", "");

    ESP.restart();
  });
  HTTP.on("/api/open", HTTP_POST, []() { httpExecute("AT+OPEN"); });
  HTTP.on("/api/close", HTTP_POST, []() { httpExecute("AT+CLOSE"); });
  HTTP.on("/api/stop", HTTP_POST, []() { httpExecute("AT+STOP"); });
  HTTP.on("/api/lock", HTTP_POST, []() { httpExecute("AT+LOCK=1"); });
  HTTP.on("/api/lock", HTTP_GET, []() { httpExecute("AT+LOCK?"); });
  HTTP.on("/api/unlock", HTTP_POST, []() { httpExecute("AT+LOCK=0"); });
  HTTP.on("/api/pos", HTTP_GET, []() { httpExecute("AT+POS?"); });
  HTTP.on("/api/pos", HTTP_POST, []() {
    String body = HTTP.arg("pos");
    if (body.length() > 3) {
      HTTP.send(400, "text/plain", "Invalid position value.");
      return;
    }
    for (byte i = 0; i < body.length(); i++) {
      if (body[i] < '0' || body[i] > '9') {
        HTTP.send(400, "text/plain", "Invalid position value.");
        return;
      }
    }

    httpExecute("AT+POS=", body.c_str());
  });

  HTTP.begin();
}

String getContentType(String filename) {
  if (filename.endsWith(".htm"))
    return "text/html";
  else if (filename.endsWith(".html"))
    return "text/html";
  else if (filename.endsWith(".css"))
    return "text/css";
  else if (filename.endsWith(".js"))
    return "application/javascript";
  else if (filename.endsWith(".png"))
    return "image/png";
  else if (filename.endsWith(".gif"))
    return "image/gif";
  else if (filename.endsWith(".jpg"))
    return "image/jpeg";
  else if (filename.endsWith(".ico"))
    return "image/x-icon";
  else if (filename.endsWith(".xml"))
    return "text/xml";
  else if (filename.endsWith(".pdf"))
    return "application/x-pdf";
  else if (filename.endsWith(".zip"))
    return "application/x-zip";
  else if (filename.endsWith(".gz"))
    return "application/x-gzip";
  return "text/plain";
}

bool handleFileRead(String path) {
  if (path.endsWith("/")) path += "index.html";
  if (!SPIFFS.exists(path)) {
    return false;
  }

  String contentType = getContentType(path);
  File file = SPIFFS.open(path, "r");
  HTTP.streamFile(file, contentType);
  file.close();
  return true;
}

void httpExecute(const char *command, const char *arg) {
  Response res = execute(command, arg);
  if (res.ok) {
    if (res.hasReturn) {
      HTTP.send(200, "text/plain", res.returnValue);
    } else {
      HTTP.send(200, "text/plain", "");
    }
  } else {
    HTTP.send(res.errCode, "text/plain", res.errDetails);
  }
}
