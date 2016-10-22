// Materia flag -- Material notification
//
// See schematics.png and physical-email-notification-video.gif
//

#include <Arduino.h>
#include <Servo.h>

// Compare numbers based on direction
class Compare
{
public:
  Compare(int step) : step_m(step) {}
  bool operator()(int a, int b) { return (step_m > 0 ? a <= b : b <= a); }
private:
  int step_m;
};

// Servo control to change angles smoothly
class SmoothServo
{
public:
  void attach(int pin)
  {
    motor_m.attach(pin);
  }

  void operator()(int new_angle, int step_length = 1, int step_delay = 5)
  {
    int step = new_angle > angle_m ? step_length : -step_length;
    Compare cmp(step);
    for(int servo_angle = angle_m; cmp(servo_angle, new_angle); servo_angle += step)
    {
      motor_m.write(servo_angle);
      delay(step_delay);
    }
    angle_m = new_angle;
    motor_m.write(angle_m);
  }
private:
  int angle_m = 0;
  Servo motor_m;
};

SmoothServo redflag;

// 9600 baud, servo attached to PIN 6
void setup() {
  redflag.attach(6);
  Serial.begin(9600);
  Serial.setTimeout(500);
}

void loop() {
  String input = Serial.readStringUntil('\n');
  if(input.length() > 0)
  {
    // valid angles are more or less in the range 0-170
    int angle = input.toInt();
    Serial.println(String("New angle: ") + String(angle));
    redflag(angle);
  }
}
