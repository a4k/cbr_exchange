@echo off
rem run this script as admin

sc create assmr-rufr-assmr.cbr.ru binpath="%CD%\server.exe" start= auto DisplayName= "assmr-rufr-assmr.cbr.ru"
sc description assmr-rufr-assmr.cbr.ru "assmr-rufr-assmr.cbr.ru"
sc start assmr-rufr-assmr.cbr.ru
sc query assmr-rufr-assmr.cbr.ru

:exit
