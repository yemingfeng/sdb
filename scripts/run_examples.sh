for((i=1;i<=10000;i++));
do
ls examples | awk '{print "go run examples/"$1}' | sh -x
done