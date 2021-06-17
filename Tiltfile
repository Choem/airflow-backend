# Load external restart proces
load('ext://restart_process', 'docker_build_with_restart')

# Infer yaml files
k8s_yaml([
  helm('deployments/file-service', name='file-service', values=['deployments/file-service/values.yaml'])
])

# Journal service deployment and live development
docker_build_with_restart(
  'k3d-airflow-backend-registry:5000/file-service', 
   './services/file-service/',
  entrypoint='/start_app',
  target='dev',
  dockerfile='./services/file-service/Dockerfile',
  live_update=[
    sync('./services/file-service/cmd', '/usr/src/app/cmd'),
  ],
)