#include <ESP8266WebServer.h>

extern ESP8266WebServer HTTP;

void initHTTPServer();
String getContentType(String filename);
bool handleFileRead(String path);
void httpExecute(const char *command, const char *arg = NULL);
