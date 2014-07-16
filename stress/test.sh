#!/bin/bash

# Note:
#		Before running this program please install ApacheBench
#		sudo apt-get install apache2-utils
if ! [ -x "$(type -P ab)" ]; then
	echo "ERROR: script requires apache bench"
	echo "For Debian and friends get it with 'apt-get install apache2-utils'"
	echo "If you have it, perhaps you don't have permissions to run it, try 'sudo $(basename $0)'"
	exit 1
fi

URL="http://109.231.121.11/api"
SESSION="YCnT0WyAHOmfRz97BNV6tnVZd7IAS17NW8dJV6OUn34sYQxIJTydL73IPR3roJDV"
HEADER_COOKIE="Cookie: DPSession=$SESSION"
HEADER_PRAGMA="Pragma: no-cache"
HEADER_ENCODING="Accept-Encoding: gzip,deflate,sdch"
HEADER_LANG="Accept-Language: en-US,en;q=0.8"
HEADER_USERAGENT="User-Agent: Mozilla/5.0 (compatible; DataPlayStress/1.0;)"
HEADER_ACCEPT="Accept: application/json, text/javascript, */*; q=0.01"
HEADER_REFERER="Referer: http://109.231.121.11/search"
HEADER_CACHE="Cache-Control: no-cache"

declare -A RUNS_LIST
RUNS_LIST=(["min"]=1 ["max"]=10)

declare -A INTENSITY_LIST
INTENSITY_LIST=(["min"]=1 ["low"]=100 ["medium"]=250 ["high"]=500 ["xhigh"]=750 ["max"]=1000)

declare -A CONCURRENCY_LIST
CONCURRENCY_LIST=(["min"]=1 ["default"]=10 ["max"]=100)

declare -A METHOD_LIST
METHOD_LIST=(["search"]="search" ["reduce"]="reduce" ["identify"]="identify" ["group"]="group")

declare -A ACTION_LIST
ACTION_LIST=(["nhs"]="nhs" ["nhsx"]="n h s" ["gold"]="gold" ["property"]="property")

RUNS=${RUNS_LIST[min]} # Number of time Run the command
INTENSITY=${INTENSITY_LIST[low]} # Number of Requests
CONCURRENCY=${CONCURRENCY_LIST[default]} # Number of concurrent requests
METHOD=${METHOD_LIST[search]}
ACTION=${ACTION_LIST[nhs]}

numRegEx='^[0-9]+$'
strRegEx='^[A-Za-z]+$'

while [[ $# > 1 ]]; do
	key="$1"
	shift

	case $key in
		-i|--intensity)
			INTENSITY="$1"
			if ! [[ $INTENSITY =~ $numRegEx ]]; then # It's not a Number
				if ! [[ $INTENSITY =~ $strRegEx ]]; then # It's neither a String
					echo "Error: intensity '$INTENSITY' is invalid!" >&2; exit 1
				elif [[ ${INTENSITY_LIST[$INTENSITY]+isset} ]]; then
					INTENSITY=${INTENSITY_LIST[$INTENSITY]}
				else # Invalid/Not in list string
					echo "Error: intensity '$INTENSITY' is not in allowed list." >&2; exit 1
				fi
			else
				if [[ ( "$INTENSITY" -lt "${INTENSITY_LIST[min]}" ) || ( "$INTENSITY" -gt "${INTENSITY_LIST[max]}" ) ]]; then # Not in Range
					echo "Error: intensity out of allowed range. (Min: ${INTENSITY_LIST[min]}, Max: ${INTENSITY_LIST[max]})" >&2; exit 1
				fi
			fi
			shift
			;;
		-c|--concurrency)
			CONCURRENCY="$1"
			if ! [[ $CONCURRENCY =~ $numRegEx ]]; then # It's not a Number
				echo "Error: concurrency '$CONCURRENCY' is invalid!" >&2; exit 1
			else
				if [[ ( "$CONCURRENCY" -lt "${CONCURRENCY_LIST[min]}" ) || ( "$CONCURRENCY" -gt "${CONCURRENCY_LIST[max]}" ) ]]; then # Not in Range
					echo "Error: concurrency out of allowed range. (Min: ${CONCURRENCY_LIST[min]}, Max: ${CONCURRENCY_LIST[max]})" >&2; exit 1
				fi
			fi
			shift
			;;
		-r|--runs)
			RUNS="$1"
			if ! [[ $RUNS =~ $numRegEx ]]; then # It's not a Number
				echo "Error: runs '$RUNS' is invalid!" >&2; exit 1
			else
				if [[ ( "$RUNS" -lt "${RUNS_LIST[min]}" ) || ( "$RUNS" -gt "${RUNS_LIST[max]}" ) ]]; then # Not in Range
					echo "Error: runs out of allowed range. (Min: ${RUNS_LIST[min]}, Max: ${RUNS_LIST[max]})" >&2; exit 1
				fi
			fi
			shift
			;;
		-m|--method)
			METHOD="$1"
			if [[ ${METHOD_LIST[$METHOD]+isset} ]]; then
				METHOD=${METHOD_LIST[$METHOD]}
			else # Invalid/Not in list string
				echo "Error: method '$METHOD' is not in allowed list." >&2; exit 1
			fi
			shift
			;;
		-a|--action)
			ACTION="$1"
			if [[ ${ACTION_LIST[$ACTION]+isset} ]]; then
				ACTION=${ACTION_LIST[$ACTION]}
			else # Invalid/Not in list string
				echo "Error: action '$ACTION' is not in allowed list." >&2; exit 1
			fi
			shift
			;;
	esac
done

site=$URL/$METHOD/$ACTION/
log=ab.$METHOD.$ACTION.log
if [ -f $log ]; then
	echo removing $log
	rm $log
fi

echo "=================================================================="
echo "= Results"
echo "=================================================================="
echo "= url ........... $site"
echo "= requests ...... $INTENSITY"
echo "= concurrency ... $CONCURRENCY"
echo "------------------------------------------------------------------"

TIMEOUT=$(( $INTENSITY*60/$CONCURRENCY ))
# echo ab -k -n $INTENSITY -c $CONCURRENCY -H "$HEADER_USERAGENT" -H "$HEADER_ACCEPT" -H "$HEADER_REFERER" -C "$COOKIE" $site
for run in $(seq 1 $runs); do
	ab -k -c $CONCURRENCY -n $INTENSITY -H "$HEADER_COOKIE" $site >> $log
	echo -e " run $run: \t $(grep "^Requests per second" $log | tail -1 | awk '{print$4}') reqs/sec"
done

avg=$(awk -v runs=$runs '/^Requests per second/ {sum+=$4; avg=sum/runs} END {print avg}' $log)

echo "------------------------------------------------------------------"
echo "= average ....... $avg requests/sec"
echo
echo "see $log for details"

# for i in {1..3}; do
# 	curl -s -w "%{time_total}\n" -o /dev/null "$URL/$METHOD/$ACTION" -H "$HEADER_PRAGMA" -H "$HEADER_ENCODING" -H "$HEADER_LANG" -H "$HEADER_USERAGENT" -H "$HEADER_ACCEPT" -H "$HEADER_REFERER" -H "$HEADER_X_REQUESTED" -H "$HEADER_COOKIE" -H "$HEADER_CONECTION" -H "$HEADER_CACHE" --compressed
# done

# End of file
