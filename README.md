## Что делает проект

Мини-контролер для умного дома. Включение-выключение света 

## Что нужно 

Плата ESP32
Кабель Micro-USB или Type-C
Макетная плата
Пара светодиодов
Провода-перемычки


## Анализ кода

```Строки: Объявление констант для работы 
    #include <WiFi.h>

    const char* ssid = "REPLACE_WITH_YOUR_SSID"; ## Объявление константы для подключение к вашему WiFi aka. Ваше имя Wifi
    const char* password = "REPLACE_WITH_YOUR_PASSWORD"; ## Пароль для вайфая 
    
    WiFiServer server(80); ## Так как наш микроконтролер будет работать на сайте, мы закидываем его на 80 порт
    String header; 
    
    String output26State = "off"; |
    String output27State = "off"; |
                                  | ## эта группа переменных для пинов, особо не шарю за верхние два со значениями "off", предполагаю что для проверки включен ли свет или нет
    const int output26 = 26;      |
    const int output27 = 27;      |
    
    unsigned long currentTime = millis(); ## Берем наше текущее время
    unsigned long previousTime = 0; 
    const long timeoutTime = 2000; ## Таймаут для выхода из программы при ошибке = 2 секунды.
```

```Строки: void setup()
        void setup() { ## Объявляем функцию для подготовки - Подключение к wifi и обработка ошибок; Назначение пинов; Подключение к серверу 
      Serial.begin(115200); 
      
      pinMode(output26, OUTPUT); | }
                                 | } ## pinMode как я понял функция которая принимает аргументом наши ПИНЫ и что они делают OUTPUT or INPUT   
      pinMode(output27, OUTPUT); | }

      digitalWrite(output26, LOW); | }
                                   | } ## digitalWrite тоже функция, берет в аргументы наши ПИНЫ и напряжение для них
      digitalWrite(output27, LOW); | }
      
      Serial.print("Connecting to "); ## Просто вывод в логи подключения к SSID
      Serial.println(ssid);
      WiFi.begin(ssid, password);
      
      while (WiFi.status() != WL_CONNECTED) { ## Цикл для повторного подключения к WiFi; 
        delay(500); 
        Serial.print(".");
      }
      
      Serial.println("\nWiFi connected."); ## Вывод успешного подключения к WiFi
      Serial.println("IP address: ");
      Serial.println(WiFi.localIP()); 
      server.begin(); ## Запуск сервера на 80 порту
    }
```

```Строки: void loop()
  void loop(){
  WiFiClient client = server.available(); ## Объявляем переменную клиента с типом WiFiClient
  
  if (client) { ## Проверка на то что клиент жив; аналог if (client == true)                          
    currentTime = millis(); ## Запускаем таймер 
    previousTime = currentTime; ## Берем самую первую секунду таймера
    Serial.println("New Client."); ## Вывод в логи       
    String currentLine = "";              
    
    while (client.connected() && currentTime - previousTime <= timeoutTime) { ## Запускаем цикл которая будет работать пока есть подключение к клиенту и пока разность наших двух переменных currentTime - previousTime меньше TimeOut
      currentTime = millis(); ## Получаем новое значение текущего временит
      if (client.available()) { ## Проверка доступности клиента            
        char c = client.read();  ## Предполагаю что это runa из моего ЯП Go, то есть один символ и мы постоянно в цикле добавляем в header символы и у нас получается какой то API запрос(GET)           
        Serial.write(c);                    
        header += c;
        
        if (c == '\n') { ## Проверка на то что дочитала строку       
          if (currentLine.length() == 0) { ## Проверка на пустую строку, значит конец запроса
            client.println("HTTP/1.1 200 OK");
            client.println("Content-type:text/html");
            client.println("Connection: close");
            client.println();
            
            if (header.indexOf("GET /26/on") >= 0) { ## Если внутри строки есть запрос ввиде "GET /26/on" - то мы включаем по 26 пину свет
              output26State = "on";
              digitalWrite(output26, HIGH);
            } else if (header.indexOf("GET /26/off") >= 0) { ## Если внутри строки есть запрос ввиде "GET /26/off" - то мы выключаем по 26 пину свет 
              output26State = "off";
              digitalWrite(output26, LOW);
            } else if (header.indexOf("GET /27/on") >= 0) {  ## Если внутри строки есть запрос ввиде "GET /26/on" - то мы включаем по 27 пину свет
              output27State = "on";
              digitalWrite(output27, HIGH);
            } else if (header.indexOf("GET /27/off") >= 0) { ## Если внутри строки есть запрос ввиде "GET /27/off" - то мы выключаем по 27 пину свет
              output27State = "off";
              digitalWrite(output27, LOW);
            }

            # Отрисовка нашего вебсайта HTML 
            client.println("<!DOCTYPE html><html>");
            client.println("<head><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">");
            client.println("<style>html { font-family: Helvetica; text-align: center;}");
            client.println(".button { background-color: #4CAF50; color: white; padding: 16px 40px; font-size: 30px; margin: 2px; cursor: pointer;}");
            client.println(".button2 {background-color: #555555;}</style></head>");
            client.println("<body><h1>ESP32 Web Server</h1>");
            
            client.println("<p>GPIO 26 - State " + output26State + "</p>");   
            if (output26State=="off") { ## Кнопка для 26 пина включение и выключение
              client.println("<p><a href=\"/26/on\"><button class=\"button\">ON</button></a></p>");
            } else {
              client.println("<p><a href=\"/26/off\"><button class=\"button button2\">OFF</button></a></p>");
            } 
               
            client.println("<p>GPIO 27 - State " + output27State + "</p>");    
            if (output27State=="off") { ## Кнопка для 27 пина для включение и выключения
              client.println("<p><a href=\"/27/on\"><button class=\"button\">ON</button></a></p>");
            } else {
              client.println("<p><a href=\"/27/off\"><button class=\"button button2\">OFF</button></a></p>");
            }
            client.println("</body></html>");
            client.println();
            break;
          } else { ## Else условие к этому if условию - if (currentLine.length() == 0). Значит то что этот else условие идет если запрос был и мы просто очищаем нашу строку для следуйщего запроса 
            currentLine = ""; 
          }
        } else if (c != '\r') { ## Проверка то что символ не != \r и добавление в следуйщую строку для запроса
          currentLine += c;      
        }
      }
    }
    header = ""; ## При остановке очищаем header
    client.stop(); ## Стопаем клиент
    Serial.println("Client disconnected.\n"); ## и выводим лог
  }
}
```



## Полный код прошивки

```cpp
#include <WiFi.h>

const char* ssid = "REPLACE_WITH_YOUR_SSID";
const char* password = "REPLACE_WITH_YOUR_PASSWORD";

WiFiServer server(80);
String header;

String output26State = "off";
String output27State = "off";
const int output26 = 26;
const int output27 = 27;

unsigned long currentTime = millis();
unsigned long previousTime = 0; 
const long timeoutTime = 2000;

void setup() {
  Serial.begin(115200);
  
  pinMode(output26, OUTPUT);
  pinMode(output27, OUTPUT);
  digitalWrite(output26, LOW);
  digitalWrite(output27, LOW);
  
  Serial.print("Connecting to ");
  Serial.println(ssid);
  WiFi.begin(ssid, password);
  
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  
  Serial.println("\nWiFi connected.");
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());
  server.begin();
}

void loop(){
  WiFiClient client = server.available();   
  
  if (client) {                             
    currentTime = millis();
    previousTime = currentTime;
    Serial.println("New Client.");          
    String currentLine = "";                
    
    while (client.connected() && currentTime - previousTime <= timeoutTime) {  
      currentTime = millis();
      if (client.available()) {             
        char c = client.read();             
        Serial.write(c);                    
        header += c;
        
        if (c == '\n') {                    
          if (currentLine.length() == 0) {
            client.println("HTTP/1.1 200 OK");
            client.println("Content-type:text/html");
            client.println("Connection: close");
            client.println();
            
            if (header.indexOf("GET /26/on") >= 0) {
              output26State = "on";
              digitalWrite(output26, HIGH);
            } else if (header.indexOf("GET /26/off") >= 0) {
              output26State = "off";
              digitalWrite(output26, LOW);
            } else if (header.indexOf("GET /27/on") >= 0) {
              output27State = "on";
              digitalWrite(output27, HIGH);
            } else if (header.indexOf("GET /27/off") >= 0) {
              output27State = "off";
              digitalWrite(output27, LOW);
            }
            
            client.println("<!DOCTYPE html><html>");
            client.println("<head><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">");
            client.println("<style>html { font-family: Helvetica; text-align: center;}");
            client.println(".button { background-color: #4CAF50; color: white; padding: 16px 40px; font-size: 30px; margin: 2px; cursor: pointer;}");
            client.println(".button2 {background-color: #555555;}</style></head>");
            client.println("<body><h1>ESP32 Web Server</h1>");
            
            client.println("<p>GPIO 26 - State " + output26State + "</p>");     
            if (output26State=="off") {
              client.println("<p><a href=\"/26/on\"><button class=\"button\">ON</button></a></p>");
            } else {
              client.println("<p><a href=\"/26/off\"><button class=\"button button2\">OFF</button></a></p>");
            } 
               
            client.println("<p>GPIO 27 - State " + output27State + "</p>");    
            if (output27State=="off") {
              client.println("<p><a href=\"/27/on\"><button class=\"button\">ON</button></a></p>");
            } else {
              client.println("<p><a href=\"/27/off\"><button class=\"button button2\">OFF</button></a></p>");
            }
            client.println("</body></html>");
            client.println();
            break;
          } else { 
            currentLine = "";
          }
        } else if (c != '\r') {  
          currentLine += c;      
        }
      }
    }
    header = "";
    client.stop();
    Serial.println("Client disconnected.\n");
  }
}
