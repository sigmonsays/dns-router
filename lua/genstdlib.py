#!/usr/bin/env python

"""
this script generates silly wrappers to stdlib functions in go so lua can call them

   % python genstdlib.py strings

"""
import sys
import os
import glob

TEMPLATE = """
package lua

import "%(package_name)s"

type %(type_name)s struct {
   %(members)s
}

func New%(type_name)s() *%(type_name)s {
   return &%(type_name)s{
      %(assign_members)s
   }
}

"""


package_name = sys.argv[1]
import_name = os.path.basename(package_name)

clean_package_name = package_name.replace("/", "_")
type_name = import_name.title()
members = set()
assign_members = set()

def match(line):
   prefix = 'pkg %s, func ' % (package_name)
   if line.startswith(prefix):
      return True, 'func'
   prefix = 'pkg %s, var ' % (package_name)
   if line.startswith(prefix):
      return True, 'var'

   return False, None

for filename in glob.glob("go*.txt"):
   for line in file(filename):
      matched, match_type = match(line)
      if matched == False:
         continue

      if match_type == 'func':
         tmp = line.split("(")
         method = tmp[0].split()[-1]
         member = "%s interface{}" % (method)
         members.add(member)

         assign = "%s: %s.%s," % (method, import_name, method)
         assign_members.add(assign)

      elif match_type == 'var':
         tmp = line.split("(")
         var_name = tmp[0].split()[-2]
         member = "%s interface{}" % (var_name)
         members.add(member)

         assign = "%s: %s.%s," % (var_name, import_name, var_name)
         assign_members.add(assign)

print "assign:", assign
print "assign_members:", assign_members


data = {
   'package_name': package_name,
   'type_name': type_name,
   'members': "\n".join(members),
   'assign_members': "\n".join(assign_members),
}


outfile = clean_package_name + '.go'
f = open(outfile, 'w+')
f.write(TEMPLATE % data)
f.close()

os.system('gofmt -w ' + outfile)
