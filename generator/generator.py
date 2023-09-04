import argparse
import yaml

def save_to_file(parsed_yaml, output):
  with open(output, "+w") as f:
    f.write(parsed_yaml)

def replace(section):
  text = section.get('text')
  field = section.get('loop').get('field')
  times = section.get('loop').get('times')

  out = ""
  for i in range(int(times)):
    out += text.replace(field, str(i+1))
  
  return out

def parse(data):
  out_yaml = ""

  for section in data.get('docker-compose'):
    if section.get('loop'):
      new = replace(section)
    else:
      new = section.get('text')
    out_yaml += new  
  return out_yaml
   
def load_yaml(filename):
  with open(filename, 'r') as stream:
    data_loaded = yaml.safe_load(stream)
  return data_loaded
   
def main():
  parser = argparse.ArgumentParser(
    prog='Docker file generator',
    description='Genera un docker file de un archivo template'
  )
  parser.add_argument('-i', '--input', help='Archivo docker-compose a armar', required=True)
  parser.add_argument('-o', '--output', help='Archivo docker-compose final', required=True)

  args = parser.parse_args()

  yaml_data = load_yaml(args.input)
  parsed_yaml = parse(yaml_data)
  save_to_file(parsed_yaml, args.output)

if __name__ == '__main__':
    main()