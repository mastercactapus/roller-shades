
struct Settings {
  int magic;
  unsigned char version;
  char name[256];      // v1
  char eventURL[256];  // v2
};

extern Settings settings;

void loadSettings();
void saveSettings();
