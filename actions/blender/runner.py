# manually assign code from this file to runnerText const in runner.go file
import sys
ind=sys.argv.index('--')+1
script=sys.argv[ind]
sys.argv=sys.argv[ind:]
exec(compile(open(script).read(), script, 'exec'), locals(), globals())
