provider "aws" {
  region = var.aws_region
}

variable "aws_region" {
  default = "us-east-1"
}

# 1. Core Networking (VPC & Subnets)
resource "aws_vpc" "sentinelx_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  tags = {
    Name = "sentinelx-soc-vpc"
  }
}

resource "aws_subnet" "public_a" {
  vpc_id            = aws_vpc.sentinelx_vpc.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "us-east-1a"
  map_public_ip_on_launch = true
}

resource "aws_subnet" "public_b" {
  vpc_id            = aws_vpc.sentinelx_vpc.id
  cidr_block        = "10.0.2.0/24"
  availability_zone = "us-east-1b"
  map_public_ip_on_launch = true
}

# 2. Database Layer (PostgreSQL)
resource "aws_db_subnet_group" "db_subnets" {
  name       = "sentinelx-db-subnets"
  subnet_ids = [aws_subnet.public_a.id, aws_subnet.public_b.id]
}

resource "aws_security_group" "db_sg" {
  name   = "sentinelx-db-security-group"
  vpc_id = aws_vpc.sentinelx_vpc.id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"] # Open only inside VPC
  }
}

resource "aws_db_instance" "postgres" {
  allocated_storage      = 20
  engine                 = "postgres"
  engine_version         = "15.4"
  instance_class         = "db.t3.micro"
  db_name                = "sentinelx"
  username               = "postgres"
  password               = "SentinelXSecurePass101!"
  db_subnet_group_name   = aws_db_subnet_group.db_subnets.name
  vpc_security_group_ids = [aws_security_group.db_sg.id]
  skip_final_snapshot    = true
}

# 3. ECS Fargate Service (API Backend)
resource "aws_ecs_cluster" "soc_cluster" {
  name = "sentinelx-soc-cluster"
}

resource "aws_security_group" "api_sg" {
  name   = "sentinelx-api-security-group"
  vpc_id = aws_vpc.sentinelx_vpc.id

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Open ingestion APIs to public agents
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_ecs_task_definition" "api_task" {
  family                   = "sentinelx-api"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"

  container_definitions = jsonencode([
    {
      name      = "api-server"
      image     = "sentinelx/server:latest"
      essential = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
        }
      ]
      environment = [
        { name = "DATABASE_URL", value = "postgres://postgres:SentinelXSecurePass101!@${aws_db_instance.postgres.endpoint}/sentinelx" },
        { name = "REDIS_URL", value = "localhost:6379" }
      ]
    }
  ])
}

resource "aws_ecs_service" "api_service" {
  name            = "sentinelx-api-service"
  cluster         = aws_ecs_cluster.soc_cluster.id
  task_definition = aws_ecs_task_definition.api_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [aws_subnet.public_a.id]
    security_groups  = [aws_security_group.api_sg.id]
    assign_public_ip = true
  }
}
