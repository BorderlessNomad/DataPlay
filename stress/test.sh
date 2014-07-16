#!/bin/bash

URL="http://109.231.121.11/api"
HEADER_ENCODING="Accept-Encoding: gzip,deflate,sdch"
HEADER_USERAGENT="User-Agent: Mozilla/5.0 (compatible; DataPlayStress/1.0;)"
HEADER_ACCEPT="Accept: application/json, text/javascript, */*; q=0.01"
HEADER_REFERER="Referer: http://109.231.121.11/search"
HEADER_CACHE="Cache-Control: no-cache"

declare -A INTENSITY_LIST
INTENSITY_LIST=(["min"]=1 ["low"]=10 ["medium"]=50 ["high"]=100 ["xhigh"]=200 ["max"]=10000)

declare -A METHOD_LIST
METHOD_LIST=(["search"]="search" ["reduce"]="reduce" ["identify"]="identify" ["group"]="group")

declare -A ACTION_LIST
ACTION_LIST=(["nhs"]="nhs" ["nhsx"]="n h s" ["gold"]="gold" ["property"]="property")

INTENSITY=${INTENSITY_LIST[low]} # Number of Requests
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

site="$URL/$METHOD/$ACTION/"
log="$METHOD.$ACTION.log"

echo "=================================================================="
echo "= Results"
echo "=================================================================="
echo "# URL ..................... $site"
echo "# Number of Requests ...... $INTENSITY"
echo "# Log file................. $log"
echo "------------------------------------------------------------------"


i="1"
total_time=0.0
while [ $i -le $INTENSITY ]; do
	run_time=$( curl -s -w "%{time_total}\n" -o /dev/null "$site" -H "$HEADER_ENCODING" -H "$HEADER_USERAGENT" -H "$HEADER_ACCEPT" -H "$HEADER_REFERER" -H "$HEADER_CACHE" )

	total_time=$(awk "BEGIN{print $total_time+$run_time}")

	i=$[$i+1]
done

avg_time=$(awk "BEGIN{print $total_time/$INTENSITY}")

echo "------------------------------------------------------------------"
echo "# Total Time (in secs)..... $total_time"
echo "# Avg. Time (in secs)...... $avg_time"
echo "------------------------------------------------------------------"

# End of file
