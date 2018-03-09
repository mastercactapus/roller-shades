#include <Arduino.h>
#include <elapsedMillis.h>

#define SERIAL_RATE 115200
#define SERIAL_TIMEOUT 1000


int dirPin = 12;
int stepPin = 11;
int runPin = 10;
int btnPin = 7;
int piPin1 = 2;
int piPin2 = 3;

bool buttonNormallyClosed = true;

long lastPos = 0;
long pos = 0;
long lastMovPos = 0;

long openStepCount = 0;
long targetSteps = -1;

elapsedMillis sinceLastMovement = 0;
elapsedMillis sinceLastState = 0;
elapsedMillis sinceLastIdle = 0;

int piState = -1;

bool ServiceLocked = true;

enum devState{
  Idle,
  IdleWait,
  Homing,
  Up,
  Down
} state;


void stateChange() {
  int newState = (digitalRead(piPin1) << 1) | digitalRead(piPin2);
  switch (piState) {
    case 0: // next is 1
      pos += newState==1 ? 1 : -1;
      break;
    case 1: // next is 3
      pos += newState==3 ? 1 : -1;
      break;
    case 3: // next is 2
      pos += newState==2 ? 1 : -1;
      break;
    case 2: // next is 0
      pos += newState==0 ? 1 : -1;
      break;
  }
  piState = newState;
  if (abs(pos - lastMovPos) > 3) {
    sinceLastMovement = 0;
    lastMovPos = pos;
  }
}

void move(unsigned long stepDelay) {
  static elapsedMicros sinceLastStep = 0;
  digitalWrite(runPin, 1);
  if (sinceLastStep > stepDelay) {
    digitalWrite(stepPin, !digitalRead(stepPin));
    sinceLastStep = 0;
  }
}

void setup() {
  // put your setup code here, to run once:
  Serial.begin(SERIAL_RATE);
  Serial.setTimeout(SERIAL_TIMEOUT);
  pinMode(btnPin, INPUT_PULLUP);
  attachInterrupt(digitalPinToInterrupt(piPin1), stateChange, CHANGE);
  attachInterrupt(digitalPinToInterrupt(piPin2), stateChange, CHANGE);
  Serial.println("Ready.");
}
void changeState(devState s) {
  if (state == s) {
    return;
  }

  state = s;
  switch (s) {
    case Up:
    case Homing:
      digitalWrite(dirPin, 1);
    case IdleWait:
      digitalWrite(runPin, 1);
      sinceLastMovement=0;
      break;
    case Idle:
      digitalWrite(runPin, 0);
      break;
    case Down:
      digitalWrite(dirPin, 0);
      digitalWrite(runPin, 1);
      sinceLastMovement=0;
      break;
  }
  sinceLastState = 0;
}

long currentPos() { // 0-100 (rounded to nearest 5)
  return (pos*200/openStepCount+5)/10*5;
}
long targetPos() { // 0-100 (rounded to nearest 5)
  return (targetSteps*200/openStepCount+5)/10*5;
}
void setTargetPos(long n) { // rounded to nearest step
  targetSteps = (openStepCount*n/10+5)/10;
}
bool isCalibrated() {
  return openStepCount != 0;
}

void STOP() {
  targetSteps = pos;
  changeState(IdleWait);
}

void processSerialCommand() {
  String c = Serial.readStringUntil('\n');
  c.trim();
  if (c.equals("")) {
    // don't respond
    return;
  } else if (c.equalsIgnoreCase("AT")) {
    // noop
  } else if (c.equalsIgnoreCase("AT+STATUS?")) {
    switch (state) {
      case Idle:
      case IdleWait:
        if (!isCalibrated()) {
          Serial.println("+STATUS:CAL");
        } else {
          Serial.println("IDLE");
        }
        break;
      case Homing:
        Serial.println("+STATUS:HOME");
        break;
      case Up:
        Serial.println("+STATUS:UP");
        break;
      case Down:
        Serial.println("+STATUS:DOWN");
        break;
    }
  } else if (c.equalsIgnoreCase("AT+OPEN")) {
    if (ServiceLocked) {
      Serial.println("ERR:700:Forbidden:Device is locked.");
      return;
    }
    setTargetPos(100);
  } else if (c.equalsIgnoreCase("AT+CLOSE")) {
    if (ServiceLocked) {
      Serial.println("ERR:700:Forbidden:Device is locked.");
      return;
    }
    setTargetPos(0);
  } else if (c.equalsIgnoreCase("AT+STOP")) {
    if (ServiceLocked) {
      Serial.println("ERR:700:Forbidden:Device is locked.");
      return;
    }
    STOP();
  } else if (c.equalsIgnoreCase("AT+LOCK?")) {
    Serial.print("+LOCK:");
    Serial.println(ServiceLocked);
  } else if (c.equalsIgnoreCase("AT+LOCK=1")) {
    ServiceLocked = true;
    STOP();
  } else if (c.equalsIgnoreCase("AT+LOCK=0")) {
    if (!isCalibrated()) {
      Serial.println("ERR:701:Not Allowed:Device must be calibrated first.");
      return;
    }
    ServiceLocked = false;
    STOP();
  } else if (c.equalsIgnoreCase("AT+POS?")) {
    if (!isCalibrated()) {
      Serial.println("ERR:501:Action Failed:Device must be calibrated first.");
      return;
    }
    Serial.print("+POS:");
    Serial.println(currentPos());
  } else if (c.substring(0,7).equalsIgnoreCase("AT+POS=")) {
    if (c.length()<8) {
      Serial.println("ERR:402:Invalid Args:New POS value is required.");
      return;
    }
    long val = c.substring(7).toInt();
    if (val < 0 || val > 100) {
      Serial.println("ERR:601:Out of Range:The argument value is not between 0 and 100 included.");
      return;
    }
    if (!isCalibrated()) {
      Serial.println("ERR:501:Action Failed:Device must be calibrated first.");
      return;
    }
    if (ServiceLocked) {
      Serial.println("ERR:700:Forbidden:Device is locked.");
      return;
    }
    long tgt = currentPos();
    if (val < tgt) {
      changeState(Down);
    } else if (val > tgt) {
      changeState(Up);
    }
    setTargetPos(val);
  } else {
    Serial.println("ERR:401:Invalid Action:Unknown action specified.");
    return;
  }

  Serial.println("OK");
}

int detectButtonPress() {
  static bool isButtonDown = false;
  static elapsedMillis sinceButtonDown = 0;
  static elapsedMillis sinceButtonUp = 0;

  int buttonDur = 0;
  if (digitalRead(btnPin) != buttonNormallyClosed) {
    if (isButtonDown && sinceButtonUp > 50) {
      buttonDur = sinceButtonDown - sinceButtonUp;
      isButtonDown = false;
    }
  } else {
    if (!isButtonDown) {
      isButtonDown = true;
      sinceButtonDown = 0;
    }
    sinceButtonUp = 0;
  }

  if (isButtonDown) {
    return 1;
  }

  return buttonDur;
}

int stepDelay(int slow, int fast, int dur) {
  return max(fast, slow-dur*5);
}

void loop() {
  if (state == Idle) {
      sinceLastIdle = 0;
  } else if (sinceLastIdle > 10000) {
    STOP();
    changeState(Idle);
    ServiceLocked = true;
  }

  int button = detectButtonPress();

  if (!button && Serial.available() >= 3) {
    processSerialCommand();
  }

  // if button is pressed and we're doing something
  // immediately stop
  if (button && state != Idle) {

    // manual stop-detection-override
    if (state == Homing) {
      openStepCount = pos;
      ServiceLocked = false;
      targetSteps = pos;
    }

    changeState(Idle);
    if (isCalibrated()) {
      targetSteps = pos;
    }
    return;
  }

  if (button >= 3000) {
    pos = 0;
    targetSteps = -1;
    openStepCount = 0;
    ServiceLocked = true;
    changeState(Homing);
  }

  long cur = currentPos();
  long tgt = targetPos();
  if (isCalibrated()) {
    if (!ServiceLocked && pos!=targetSteps && state == Idle) {
      if (cur >= 95) { // pulled down from top
        setTargetPos(0);
        tgt = 0;
      } else if (cur < 95) {
        setTargetPos(100);
        tgt = 100;
      }
    }
    if (cur > tgt) {
      changeState(Down);
    } else if (cur < tgt) {
      changeState(Up);
    } else if (state == Down && pos <= targetSteps) {
      changeState(IdleWait);
    } else if (state == Up && pos >= targetSteps) {
      changeState(IdleWait);
    }
  }

  switch (state) {
    case Idle:
      // nothing to do
      break;
    case Homing:
      if (sinceLastMovement > 1000) {
        openStepCount = pos;
        ServiceLocked = false;
        targetSteps = pos;
        changeState(IdleWait);
        return;
      }
      move(2000);
      break;
    case Up:
      if (sinceLastMovement > 600) {
        STOP();
        return;
      }

      move(stepDelay(2200, 1300, sinceLastState));
      break;
    case Down:
      if (sinceLastMovement > 600) {
        STOP();
        return;
      }

      move(stepDelay(2200, 750, sinceLastState));
      break;
    case IdleWait:
      if (sinceLastState > 1500) {
        targetSteps = pos;
        changeState(Idle);
        return;
      }
      if (sinceLastState > 1000) {
        digitalWrite(runPin, 0);
      } else {
        digitalWrite(runPin, 1);
      }
  }

}
