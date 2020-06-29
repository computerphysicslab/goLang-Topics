# goLang-Topics
Get github topics from repositories matching one or several topics

This goLang code finds tag/topic recommendations by using a real-time algorithm that searches github repositories matching a given term, and extracts their topics. Most frequent topics are selected as the result.

For example, exploring relevant topics related to "nlp" the following recommendations are found, along with their frequencies:

RESULT: [{nlp 30} {natural-language-processing 18} {machine-learning 16} {deep-learning 11} {tensorflow 10} {python 9} {data-science 5} {neural-network 5} {pytorch 5} {bert 4}]

Or when exploring results for lda:

RESULT: [{lda 30} {topic-modeling 13} {machine-learning 10} {nlp 9} {pca 7} {latent-dirichlet-allocation 6} {gibbs-sampling 5} {natural-language-processing 5} {topic-models 3} {python 3}]

An so on... 

Results for go-chi: [{go-chi 30} {golang 19} {go 10} {rest-api 5} {rest 4} {mongodb 4} {middleware 4} {mysql 4} {docker 4} {web 3}]

Results for gitea: [{gitea 30} {github 9} {gitlab 9} {git 9} {docker 8} {gogs 7} {bitbucket 4} {traefik 4} {golang 4} {ansible 4}]

It is posible to query two topics by adding a + between them:

Results for ansible+traefik: [{ansible 30} {traefik 27} {docker 15} {docker-compose 8} {ansible-role 7} {kubernetes 5} {ansible-playbook 5} {docker-swarm 4} {terraform 4} {devops 4}]
