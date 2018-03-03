#include <elapsedMillis.h>

int dirPin = 12;
int stepPin = 11;
int runPin = 10;
int btnPin = 7;
int piPin1 = 2;
int piPin2 = 3;

int lastPos = 0;
int pos = 0;
int lastMovPos = 0;

int upPos = 0;
int downPos = 0;

bool hasUpPos = false;
bool hasDownPos = false;
bool isButtonDown = false;
bool goingDown = false;

elapsedMicros sinceLastStep;
elapsedMillis sinceLastMovement;
elapsedMillis sinceButtonDown;
elapsedMillis sinceButtonUp;

int piState = -1;

enum devState{
  Idle,
  IdleWait,
  Homing,
  Up,
  Down,
  Recover
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
  digitalWrite(runPin, 1);
  if (sinceLastStep > stepDelay) {
    digitalWrite(stepPin, !digitalRead(stepPin));
    sinceLastStep = 0;
  }
}

void setup() {
  // put your setup code here, to run once:
  Serial.begin(115200);
  pinMode(btnPin, INPUT_PULLUP);
  attachInterrupt(digitalPinToInterrupt(piPin1), stateChange, CHANGE);
  attachInterrupt(digitalPinToInterrupt(piPin2), stateChange, CHANGE);
}
void changeState(devState s) {
  if (state == s) {
    return;
  }
  Serial.print("State=");
  switch (s) {
    case Idle:
    Serial.print("Idle");
    break;
    case IdleWait:
    Serial.print("IdleWait");
    break;
    case Up:
    Serial.print("Up");
    break;
    case Down:
    Serial.print("Down");
    break;
    case Homing:
    Serial.print("Homing");
    break;
    case Recover:
    Serial.print("Recover");
    break;
  }
  Serial.print(",Pos=");
  Serial.print(pos);
  Serial.print(",Up=");
  Serial.print(upPos);
  Serial.print(",Dn=");
  Serial.println(downPos);

  state = s;
  switch (s) {
    case Up:
    case Recover:
    case Homing:
      goingDown = false;
      digitalWrite(dirPin, 1);
    case IdleWait:
      digitalWrite(runPin, 1);
      sinceLastMovement=0;
      break;
    case Idle:
      digitalWrite(runPin, 0);
      break;
    case Down:
      goingDown = true;
      digitalWrite(dirPin, 0);
      digitalWrite(runPin, 1);
      sinceLastMovement=0;
      break;
  }
}
bool isTop() {
  return pos > (upPos-10);
}
void toggle() {
  sinceLastMovement=0;
  // short press should home if needed
  // otherwise toggle if possible
  if (!hasUpPos) {
    changeState(Homing);
  } else if ((isTop() || !goingDown) && hasDownPos) {
    changeState(Down);
  } else if (!isTop() && hasUpPos && goingDown) {
    changeState(Up);
  }
}

void loop() {
  int buttonDur = 0;

  if (digitalRead(btnPin)) {
    if (isButtonDown && sinceButtonUp > 100) {
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


  if (!isButtonDown && hasDownPos && pos < downPos-10 && state != Up) {
    changeState(Recover);
  }
  
  if (buttonDur > 0) {
    Serial.print("Dur=");
    Serial.println(buttonDur);
    isButtonDown = false;
  }

  // if button is pressed and we're doing something
  // immediately stop
  if (isButtonDown && (state==Homing||state==Up||state==Down||state==Recover)) {
    changeState(IdleWait);
    return;
  }

  if (buttonDur < 35) {
    // do nothing for "blips"
  } else if (buttonDur < 1200 && state == Idle) {
    toggle();
  } else if (buttonDur >= 1200) {
    downPos = pos;
    hasDownPos = true;
    hasUpPos = false;
    Serial.print("Down=");
    Serial.println(downPos);
  }

  switch (state) {
    case Idle:
      if (!isButtonDown && hasDownPos && hasUpPos && sinceLastMovement < 100) {
        toggle();
        return;
      }
      break;
    case Homing:
      if (sinceLastMovement > 1000) {
        upPos = pos;
        hasUpPos = true;
        Serial.print("Up=");
        Serial.println(upPos);
        changeState(IdleWait);
        return;
      }
      move(2000);
      break;
    case Recover:
      if (pos >= downPos+10 || sinceLastMovement > 600) {
        changeState(IdleWait);
        return;
      }
      move(2000);
      break;
    case Up:
      if (pos >= upPos || sinceLastMovement > 600) {
        changeState(IdleWait);
        return;
      }
      if (sinceLastMovement < 100) {
        move(1500);
      } else {
        move(2000);
      }
      break;
    case Down:
      if (pos <= downPos || sinceLastMovement > 600) {
        changeState(IdleWait);
        return;
      }
      move(750);
      break;
    case IdleWait:
      if (sinceLastMovement > 600) {
        changeState(Idle);
        return;
      }
      if (sinceLastMovement > 300) {
        digitalWrite(runPin, 0);
      } else {
        digitalWrite(runPin, 1);
      }
  }

}

