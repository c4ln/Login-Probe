# Login-Probe
A Go-based tool for educational purposes that explores potential security vulnerabilities in IoT devices through HTTP authentication brute force testing.

This project is meant to demonstrate how inadequate security practices in IoT devices can expose them to unauthorized access.

## Disclaimer

This project is intended solely for educational purposes. The goal is to raise awareness about the importance of securing IoT devices and to encourage ethical and responsible cybersecurity practices. The authors of this project do not condone or endorse any illegal or malicious activities performed with the information and code provided here. Use this tool responsibly and in compliance with all applicable laws and regulations.

## Features

- HTTP authentication brute force testing against IoT devices.
- Identification of potential vulnerabilities in IoT device authentication.
- Educational insights into the security risks associated with IoT devices.

## Usage

1. Ensure you have Go installed on your system.
2. Clone this repository to your local machine.
3. Edit the `logins.txt` file to include the list of logins you want to test. Ensure the logins.txt are in the form of user:password
4. Run the program using the following command:

   ```shell
   ./main <port/listen> <http/https> <exploit 1=yes,0=no>

## Usage example:
zmap -p (port) -q | ./main (port used in zmap) http 0

Ensure, the port of the IP is the same port given as the first argument. 
