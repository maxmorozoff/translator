#!/bin/bash
PRJNAME=translate
processes=$(ps aux | grep $PRJNAME -i|grep -v grep | awk -F ' ' '{print $2}' | xargs)
echo 'Killing processes... '$processes
for i in $processes; do sudo kill $i; done