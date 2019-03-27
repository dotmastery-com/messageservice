node {
    
    // Install the desired Go version
    echo "Installing go"
    def root = tool name: '1.12.1', type: 'go'
    echo "done"
    
    sh 'pwd'
    
    try{
       
        
        ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/realtime-chat") {
            withEnv(["GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}", "GOROOT=${root}","PATH+GO=${root}/bin"]) {
                env.PATH="${GOPATH}/bin:$PATH"
                
                stage('Checkout'){
                    echo 'Checking out SCM'
                    checkout([$class: 'GitSCM', branches: [[name: '*/master']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[url: 'https://github.com/dotmastery-com/messageservice.git']]])
                }
                
                stage('Pre Test'){
                    echo 'Pulling Dependencies'
            
                    sh 'go version'
                    sh 'go get github.com/satori/go.uuid'
                    sh 'go get github.com/gorilla/websocket'
                    sh 'go get github.com/rs/xid'
                    sh 'go get github.com/segmentio/kafka-go'

                }
        
                stage('Test'){
                    
                    //List all our project files with 'go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org'
                    //Push our project files relative to ./src
                   // sh 'cd $GOPATH && go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org > projectPaths'
                    
                    //Print them with 'awk '$0="./src/"$0' projectPaths' in order to get full relative path to $GOPATH
                   // def paths = sh returnStdout: true, script: """awk '\$0="./src/"\$0' projectPaths"""
                    
                  //  echo 'Vetting'

//                    sh """cd $GOPATH && go tool vet ${paths}"""

                  //  echo 'Linting'
  //                  sh """cd $GOPATH && golint ${paths}"""
                    
               //     echo 'Testing'
    //                sh """cd $GOPATH && go test -race -cover ${paths}"""
                }
            
                stage('Build'){
                    echo 'Building Executable'
                    
                    sh 'pwd'
                    sh 'tree'
                
                    //Produced binary is $GOPATH/src/cmd/project/project
                    sh """go build -ldflags '-s'"""
                }
                
                
                
            }
        }
    }catch (e) {
        // If there was an exception thrown, the build failed
        currentBuild.result = "FAILED"
        
        
    } finally {
        
    }
}