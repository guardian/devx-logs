FROM amazon/aws-for-fluent-bit:init-2.32.0

ENV aws_fluent_bit_init_file_1 /guardian.conf
ENV aws_fluent_bit_init_file_2 /guardian-parsers.conf

# Added by aws-for-fluent-bit base image
# https://github.com/aws/aws-for-fluent-bit/blob/develop/scripts/dockerfiles/Dockerfile.main-release#L49
ENV aws_fluent_bit_init_file_3 /ecs/ecs-metadata.conf

ADD guardian-parsers.conf /guardian-parsers.conf
ADD guardian.conf /guardian.conf
