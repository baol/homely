// Materia flag -- Material notification
//
// See ../schematics.png and ../physical-email-notification-flag-video.gif
//

#include <Arduino.h>
#include <Servo.h>

// Servo control to change angles smoothly
class SmoothServo
{
public:

  void attach(int pin)
  {
    motor_m.attach(pin);
  }

  void operator()(int new_angle, int duration = 750, int total_steps = 30)
  {
    Interpolate interpolator(angle_m, new_angle, total_steps);
    for(int ii(0); ii!=total_steps; ++ii)
    {
      motor_m.write(interpolator.pos_at_step(ii));
      delay(duration/total_steps);
    }
    angle_m = new_angle;
    motor_m.write(angle_m);
  }

private:

  class Interpolate
  {
  public:
    Interpolate(int from , int to, int steps)
      : m_from(from)
      , m_to(to)
      , m_steps(steps)
    {
    }

    int pos_at_step(int step)
    {
      double pos = 0;
      if(step < m_steps/2)
      {
        // first parabolic curve: acceleration
        double sqrt_pos = 1.414*double(step)/m_steps;
        pos = sqrt_pos * sqrt_pos;
      }
      else
      {
        // second parabolic curve for deceleration
        double sqrt_pos = 1.414*double(step)/m_steps - 1.414;
        pos = -sqrt_pos * sqrt_pos + 1;
      }
      return m_from + int(double(m_to-m_from)*pos);
    }

  private:
    int m_from, m_to, m_steps;
  };

  int angle_m = 90;
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
    redflag(angle, 900, 30);
  }
}
