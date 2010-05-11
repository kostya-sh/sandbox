"""
  Usage: countline.py sourcedir resultfile
"""
import sys
import os
import os.path

VALID_EXTS = set(['.py', '.conf', '.xml', '.vm', '.bat', '.js', '.m', 
   '.sql', '.jnlp', '.css', '.cfg', '.txt', '.html', '.java'])

def count_lines(f):
  return len(filter(lambda s: len(s.strip()) != 0, open(f, 'r')))

def count_lines_by_ext(rootdir):
    counts = {} # ext -> count
    for root, dirs, files in os.walk(rootdir):
        if '.svn' in dirs:
            dirs.remove('.svn')
        for f in files:
            fn, fext = os.path.splitext(f.lower())
            if fext == ".htm":
                fext = ".html"
            if fext in VALID_EXTS:
                cnt = counts.get(fext, 0)
                cnt += count_lines(os.path.join(root, f))
                counts[fext] = cnt
    counts["total"] = sum(counts.values())

    return counts

res = "\n".join(map(lambda s: "%s: %d" % (s[0], s[1]), count_lines_by_ext(sys.argv[1]).items()))
try:
    f = open(sys.argv[2], 'w')
    f.write(res)
except:
    print res
