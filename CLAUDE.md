# Notiflex Platform

## 프로젝트 개요

Notiflex — B2B 알림 SaaS 플랫폼. 이 저장소는 애플리케이션 코드, Kubernetes 매니페스트, CI/CD 파이프라인을 관리한다.

## 기술 스택

- **언어**: Go 표준 라이브러리 (외부 프레임워크 없음)
- **컨테이너**: scratch 베이스 이미지

## GCP 설정

- **프로젝트 ID**: project-d64f9b5c-20c8-4906-95b
- **리전**: asia-northeast3 (서울)
- **존**: asia-northeast3-a

## Artifact Registry

```
asia-northeast3-docker.pkg.dev/project-d64f9b5c-20c8-4906-95b/notiflex
```

## 행동 규칙

- 클러스터/인프라 변경 전 항상 현재 상태를 확인한다.
- 파괴적이거나 되돌리기 어려운 작업(삭제, force-push 등)은 실행 전 사용자에게 확인받는다.
- 모든 kubectl 명령에는 `--context gke_project-d64f9b5c-20c8-4906-95b_asia-northeast3-a_notiflex-cluster`를 명시한다.
