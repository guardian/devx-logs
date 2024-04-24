FROM amazon/aws-for-fluent-bit:init-2.32.0

# Environment variables are used to specify which config files to load
# https://github.com/aws/aws-for-fluent-bit/blob/mainline/use_cases/init-process-for-fluent-bit/README.md#how-to-use-multi-config-feature

ENV aws_fluent_bit_init_file_1 /guardian-parsers.conf

# Added by aws-for-fluent-bit base image
# https://github.com/aws/aws-for-fluent-bit/blob/develop/scripts/dockerfiles/Dockerfile.main-release#L49
ENV aws_fluent_bit_init_file_2 /ecs/ecs-metadata.conf

# This config file must go last as it applies a `parser` filter which will interfere with any filters that are run before it (config files loaded in reverse order)
ENV aws_fluent_bit_init_file_3 /guardian.conf

ADD guardian-parsers.conf /guardian-parsers.conf
ADD guardian.conf /guardian.conf

VOLUME [ "/init" ]