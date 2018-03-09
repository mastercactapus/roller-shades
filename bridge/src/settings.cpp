#include <EEPROM.h>

#include "settings.h"

Settings settings;

void saveSettings() {
  EEPROM.put(0, settings);
  EEPROM.commit();
}
void loadSettings() {
  EEPROM.begin(1024);
  EEPROM.get(0, settings);

  // migrate settings with defaults as "versions" are added
  if (settings.magic != 0x1337) {
    settings.magic = 0x1337;
    settings.version = 0;
  }

  if (settings.version < 1) {
    settings.magic = 0x1337;
    settings.version = 1;
    strcpy(settings.name, "DIY Roller Shades");
  }

  if (settings.version < 2) {
    settings.eventURL[0] = 0;
    settings.version = 2;
  }

  saveSettings();
}
