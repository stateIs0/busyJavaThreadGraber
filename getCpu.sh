#!/bin/sh
source /etc/profile

#define variable
psUser=$1
psProcess=$2
pid= `ps -ef | egrep ${psProcess} | egrep ${psUser} |  egrep -v "grep|vi|tail" | sed -n 1p | awk '{print $2}'`
echo ${pid}
if [ -z ${pid} ];then
	echo "The process does not exist."
	exit 1
fi


CpuValue=`ps -p ${pid} -o pcpu |egrep -v CPU |awk '{print $1}'`

echo ${CpuValue}

flag=`echo ${CpuValue} | awk -v tem=80 '{print($1>tem)? "1":"0"}'`

if [ ${flag} -eq 1 ];then

        echo “The usage of cpu is larger than 80%”
else
        echo “The usage of cpu is ok”
fi

