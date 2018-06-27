# TL;dr

It's create-react-app for Dockerfiles

# A Dev Workflow using `app`

1. To start off, just write your app. 
1. Now, to run it, do `app run .`. This runs your app in a container. By default it will create you an Appfile based on the language it detects (if it can't, it will error).
1. To run the tests, do `app run . -c scripts/test` (it's nice to not use opaque container images!)
1. The container is content-addressed, rebuildable, and reproducible - using the git hash and the docker image hash
1. Let's say you would like to run someone else's code, from github.. push the code to github now
1. `app run github.com/williammartin/myapp`. That's all you need to do to run it.
1. What about pushing the app to CF, to a CI system, or to Kubernetes?
1. `app build github.com/williammartin/myapp -t foo/bar; docker push foo/bar; cf push -o foo/bar`
1. Bonus points: `cf push github.com/williammartin/myapp` (cf responsible for patching your rootfs)
1. Bonus points: `app ci .` builds concourse yml.
