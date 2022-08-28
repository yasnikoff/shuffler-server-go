package blender
const runnerText=`
import sys
ind=sys.argv.index('--')+1
script=sys.argv[ind]
sys.argv=sys.argv[ind:]
exec(compile(open(script).read(), script, 'exec'), locals(), globals())
`
