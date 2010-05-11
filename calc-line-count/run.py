import sys
import os
import subprocess
from datetime import date, timedelta

checkout_base = sys.argv[1]

one_day = timedelta(days = 1)
t = max([int(s[:8]) for s in os.listdir("results")])
initial_date = date(t / 10000, (t % 10000) / 100, t % 100) + one_day
last_date = date.today()

d = initial_date
while d <= last_date:
    print d
    svn_date = "{%d-%02d-%02d 00:00:00}" % (d.year, d.month, d.day)
    res_file = "results/%d%02d%02d.txt" % (d.year, d.month, d.day)

    subprocess.call(["svn", "update", "-r", svn_date, checkout_base])
    print "Calculating line count...",
    subprocess.call(["python", "countlines.py", checkout_base, res_file])
    print "DONE: " + res_file
    
    d += one_day
    print 80 * "=", "\n"

print "Visualising...",
subprocess.call(["python", "visualize.py", "loc.png"])
print "DONE"
