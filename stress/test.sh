#!/bin/bash

URL="http://109.231.121.11/api/"

declare -A INTENSITY_LIST
INTENSITY_LIST=(["min"]=1 ["low"]=100 ["medium"]=250 ["high"]=500 ["xhigh"]=750 ["max"]=1000)

declare -A FREQUENCY_LIST
FREQUENCY_LIST=(["min"]=1 ["max"]=60)

declare -A METHOD_LIST
METHOD_LIST=(["search"]="search" ["reduce"]="reduce" ["identify"]="identify" ["group"]="group")

declare -A ACTION_LIST
ACTION_LIST=(["nhs"]="nhs" ["nhsx"]="n h s" ["gold"]="gold" ["property"]="property")

INTENSITY=${INTENSITY_LIST[low]} # Number of Requests
FREQUENCY=${FREQUENCY_LIST[min]} # Per N Seconds
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
		-f|--frequency)
			FREQUENCY="$1"
			if ! [[ $FREQUENCY =~ $numRegEx ]]; then # It's not a Number
				echo "Error: frequency '$FREQUENCY' is invalid!" >&2; exit 1
			else
				if [[ ( "$FREQUENCY" -lt "${FREQUENCY_LIST[min]}" ) || ( "$FREQUENCY" -gt "${FREQUENCY_LIST[max]}" ) ]]; then # Not in Range
					echo "Error: frequency out of allowed range. (Min: ${FREQUENCY_LIST[min]}, Max: ${FREQUENCY_LIST[max]})" >&2; exit 1
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
		--default)
			DEFAULT=YES
			shift
			;;
		*)
			error "Unexpected option $1"
			;;
	esac
done

echo "intensity = ${INTENSITY}"
echo "frequency = ${FREQUENCY}"
echo "method = ${METHOD}"
echo "action = ${ACTION}"

# End of file
