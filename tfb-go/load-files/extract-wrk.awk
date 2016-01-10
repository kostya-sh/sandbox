/connections/ { T=$1; C=$4; if (T == "1") {print("Threads\tConns\tReqs/s\tSpeed\tAvg Lat\tMax Lat\t")} }
/Latency/ { LA=$2; LM=$4 }
/Requests/ { R=$2 }
/Transfer/ { S=$2 ; print(T "\t" C "\t" R "\t" S "\t" LA "\t" LM) }
