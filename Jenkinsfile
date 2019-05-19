pipeline {
    agent none



    environment {
        registry = "dotmastery/messageservice"
    }
    stages {


        stage('Build and Push Docker Image') {
          agent any  

          steps {
            script {
                
                docker.withRegistry('https://registry-1.docker.io/v2/', 'Dockerhub') {
                    dockerImage = docker.build registry
                    dockerImage.push()
                }    
            }    

         }
        }

         stage('Publish to Openshift') {
            agent any 

             steps {
                 script {
                    sh '/app/openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit/oc login https://api.pro-eu-west-1.openshift.com --token=eT-2btD45f6QeFWhgLNMj3GPqgC5rA0SAD7ZbNz5EMU'
                    sh '/app/openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit/oc import-image messageservice:latest --from=dotmastery/messageservice --confirm'
                }
             }  
         }

    }
}