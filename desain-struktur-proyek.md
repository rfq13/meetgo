# Desain Struktur Folder Proyek untuk Aplikasi WebRTC Meeting

## 1. Pendahuluan

Struktur folder proyek ini dirancang untuk mendukung pengembangan aplikasi WebRTC meeting dengan backend Golang dan frontend Preact. Struktur ini mengikuti best practices untuk organisasi kode yang maintainable, scalable, dan collaborative.

## 2. Struktur Folder Root

```
webrtc-meeting-app/
├── .github/                     # GitHub workflows dan issue templates
├── backend/                     # Backend application (Golang)
├── frontend/                    # Frontend application (Preact)
├── deployments/                 # Deployment configurations
├── docs/                        # Documentation
├── scripts/                     # Utility scripts
├── .env.example                 # Environment variables example
├── .gitignore                   # Git ignore rules
├── docker-compose.yml           # Docker Compose for local development
├── Dockerfile                   # Dockerfile for the application
├── LICENSE                      # License file
├── Makefile                     # Build automation
├── package.json                 # Root package.json for shared scripts
├── README.md                    # Project README
└── .vscode/                     # VSCode configuration
    ├── extensions.json          # Recommended extensions
    ├── launch.json              # Debug configurations
    └── settings.json             # VSCode settings
```

## 3. Backend Structure (Golang)

```
backend/
├── cmd/                         # Application entry points
│   ├── api/                     # API server
│   │   └── main.go              # API server entry point
│   ├── websocket/               # WebSocket server
│   │   └── main.go              # WebSocket server entry point
│   ├── media/                   # Media server
│   │   └── main.go              # Media server entry point
│   └── migration/               # Database migration tool
│       └── main.go              # Migration tool entry point
├── internal/                    # Private application code
│   ├── api/                     # HTTP handlers and middleware
│   │   ├── middleware/           # API middleware
│   │   │   ├── auth.go          # Authentication middleware
│   │   │   ├── cors.go          # CORS middleware
│   │   │   ├── logging.go       # Logging middleware
│   │   │   ├── rate_limit.go    # Rate limiting middleware
│   │   │   └── recovery.go      # Recovery middleware
│   │   ├── v1/                   # API version 1
│   │   │   ├── auth/            # Authentication handlers
│   │   │   │   ├── handler.go   # Auth handlers
│   │   │   │   └── routes.go    # Auth routes
│   │   │   ├── users/           # User handlers
│   │   │   │   ├── handler.go   # User handlers
│   │   │   │   └── routes.go    # User routes
│   │   │   ├── rooms/           # Room handlers
│   │   │   │   ├── handler.go   # Room handlers
│   │   │   │   └── routes.go    # Room routes
│   │   │   ├── notifications/   # Notification handlers
│   │   │   │   ├── handler.go   # Notification handlers
│   │   │   │   └── routes.go    # Notification routes
│   │   │   ├── webrtc/          # WebRTC handlers
│   │   │   │   ├── handler.go   # WebRTC handlers
│   │   │   │   └── routes.go    # WebRTC routes
│   │   │   └── admin/           # Admin handlers
│   │   │       ├── handler.go   # Admin handlers
│   │   │       └── routes.go    # Admin routes
│   │   └── router.go            # API router
│   ├── auth/                    # Authentication service
│   │   ├── handler.go            # Auth handler interface
│   │   ├── service.go            # Auth service implementation
│   │   ├── repository.go         # Auth repository interface
│   │   └── repository_impl.go    # Auth repository implementation
│   ├── config/                   # Configuration management
│   │   ├── config.go             # Config struct
│   │   ├── config_test.go        # Config tests
│   │   └── loader.go             # Config loader
│   ├── database/                 # Database connection and migrations
│   │   ├── connection.go         # Database connection
│   │   ├── migrations/           # Database migrations
│   │   │   ├── 000001_create_users_table.up.sql
│   │   │   ├── 000001_create_users_table.down.sql
│   │   │   ├── 000002_create_rooms_table.up.sql
│   │   │   ├── 000002_create_rooms_table.down.sql
│   │   │   ├── 000003_create_room_participants_table.up.sql
│   │   │   ├── 000003_create_room_participants_table.down.sql
│   │   │   ├── 000004_create_user_contacts_table.up.sql
│   │   │   ├── 000004_create_user_contacts_table.down.sql
│   │   │   ├── 000005_create_user_sessions_table.up.sql
│   │   │   ├── 000005_create_user_sessions_table.down.sql
│   │   │   ├── 000006_create_room_messages_table.up.sql
│   │   │   ├── 000006_create_room_messages_table.down.sql
│   │   │   ├── 000007_create_meeting_history_table.up.sql
│   │   │   ├── 000007_create_meeting_history_table.down.sql
│   │   │   ├── 000008_create_user_settings_table.up.sql
│   │   │   ├── 000008_create_user_settings_table.down.sql
│   │   │   ├── 000009_create_room_settings_table.up.sql
│   │   │   ├── 000009_create_room_settings_table.down.sql
│   │   │   ├── 000010_create_notifications_table.up.sql
│   │   │   ├── 000010_create_notifications_table.down.sql
│   │   │   ├── 000011_create_functions_and_triggers.up.sql
│   │   │   ├── 000011_create_functions_and_triggers.down.sql
│   │   │   ├── 000012_create_views.up.sql
│   │   │   ├── 000012_create_views.down.sql
│   │   │   ├── 000013_setup_security.up.sql
│   │   │   └── 000013_setup_security.down.sql
│   │   └── seeds/                 # Database seeds
│   │       ├── users.sql          # User seeds
│   │       └── rooms.sql          # Room seeds
│   ├── models/                   # Database models
│   │   ├── user.go               # User model
│   │   ├── room.go               # Room model
│   │   ├── room_participant.go   # Room participant model
│   │   ├── user_contact.go       # User contact model
│   │   ├── user_session.go       # User session model
│   │   ├── room_message.go       # Room message model
│   │   ├── meeting_history.go    # Meeting history model
│   │   ├── user_setting.go       # User setting model
│   │   ├── room_setting.go       # Room setting model
│   │   └── notification.go       # Notification model
│   ├── notification/             # Notification service
│   │   ├── handler.go            # Notification handler interface
│   │   ├── service.go            # Notification service implementation
│   │   ├── repository.go         # Notification repository interface
│   │   ├── repository_impl.go    # Notification repository implementation
│   │   ├── email.go              # Email notification
│   │   ├── push.go               # Push notification
│   │   └── webhook.go            # Webhook notification
│   ├── room/                     # Room service
│   │   ├── handler.go            # Room handler interface
│   │   ├── service.go            # Room service implementation
│   │   ├── repository.go         # Room repository interface
│   │   └── repository_impl.go    # Room repository implementation
│   ├── user/                     # User service
│   │   ├── handler.go            # User handler interface
│   │   ├── service.go            # User service implementation
│   │   ├── repository.go         # User repository interface
│   │   └── repository_impl.go    # User repository implementation
│   ├── webrtc/                   # WebRTC service
│   │   ├── handler.go            # WebRTC handler interface
│   │   ├── peer.go               # WebRTC peer connection
│   │   ├── room.go               # WebRTC room
│   │   ├── sfu.go                # SFU implementation
│   │   └── media.go              # Media handling
│   ├── websocket/                # WebSocket service
│   │   ├── hub.go                # WebSocket hub
│   │   ├── client.go             # WebSocket client
│   │   ├── handler.go            # WebSocket handlers
│   │   └── types.go              # WebSocket types
│   └── utils/                    # Utility functions
│       ├── crypto.go             # Cryptographic utilities
│       ├── validator.go          # Validation utilities
│       ├── response.go           # Response utilities
│       ├── time.go               # Time utilities
│       └── string.go             # String utilities
├── pkg/                         # Public library code
│   ├── logger/                   # Logger implementation
│   │   ├── logger.go             # Logger interface
│   │   ├── logrus.go            # Logrus implementation
│   │   └── zap.go               # Zap implementation
│   └── metrics/                  # Metrics implementation
│       ├── metrics.go            # Metrics interface
│       ├── prometheus.go         # Prometheus implementation
│       └── middleware.go         # Metrics middleware
├── scripts/                     # Scripts
│   ├── migration.sh              # Migration script
│   ├── build.sh                  # Build script
│   ├── test.sh                   # Test script
│   ├── deploy.sh                 # Deployment script
│   └── docker-entrypoint.sh      # Docker entrypoint script
├── configs/                     # Configuration files
│   ├── app.yaml                  # Application config
│   ├── database.yaml             # Database config
│   ├── redis.yaml                # Redis config
│   ├── webrtc.yaml               # WebRTC config
│   └── logging.yaml              # Logging config
├── deployments/                 # Deployment configurations
│   ├── docker/                   # Docker configs
│   │   ├── Dockerfile.api        # Dockerfile for API server
│   │   ├── Dockerfile.websocket  # Dockerfile for WebSocket server
│   │   ├── Dockerfile.media      # Dockerfile for Media server
│   │   └── docker-compose.yml    # Docker Compose for local development
│   └── kubernetes/               # Kubernetes configs
│       ├── api-deployment.yaml   # API deployment
│       ├── websocket-deployment.yaml # WebSocket deployment
│       ├── media-deployment.yaml # Media deployment
│       ├── postgres-deployment.yaml # PostgreSQL deployment
│       ├── redis-deployment.yaml # Redis deployment
│       ├── api-service.yaml      # API service
│       ├── websocket-service.yaml # WebSocket service
│       ├── media-service.yaml    # Media service
│       ├── configmap.yaml        # ConfigMap
│       └── secret.yaml           # Secret
├── tests/                       # Test files
│   ├── integration/              # Integration tests
│   │   ├── auth_test.go         # Auth integration tests
│   │   ├── user_test.go         # User integration tests
│   │   ├── room_test.go         # Room integration tests
│   │   └── webrtc_test.go       # WebRTC integration tests
│   ├── unit/                     # Unit tests
│   │   ├── auth_test.go         # Auth unit tests
│   │   ├── user_test.go         # User unit tests
│   │   ├── room_test.go         # Room unit tests
│   │   ├── websocket_test.go    # WebSocket unit tests
│   │   └── webrtc_test.go       # WebRTC unit tests
│   └── e2e/                      # End-to-end tests
│       ├── auth_test.go         # Auth e2e tests
│       ├── meeting_test.go      # Meeting e2e tests
│       └── admin_test.go        # Admin e2e tests
├── go.mod                       # Go modules
├── go.sum                       # Go modules checksum
├── golangci.yml                 # GolangCI-Lint configuration
├── .air.toml                    # Air live reload configuration
├── .golangci.yml                # GolangCI-Lint configuration
├── .goreleaser.yml              # GoReleaser configuration
├── Makefile                     # Build automation
└── README.md                     # Backend README
```

## 4. Frontend Structure (Preact)

```
frontend/
├── public/                      # Static assets
│   ├── favicon.ico               # Favicon
│   ├── manifest.json             # PWA manifest
│   ├── robots.txt                # Robots.txt
│   └── assets/                   # Images, fonts, etc.
│       ├── icons/                 # SVG icons
│       │   ├── audio.svg
│       │   ├── video.svg
│       │   ├── microphone.svg
│       │   ├── videocam.svg
│       │   ├── screen-share.svg
│       │   ├── chat.svg
│       │   ├── people.svg
│       │   ├── settings.svg
│       │   └── logout.svg
│       ├── images/                # Images
│       │   ├── logo.png
│       │   ├── hero-bg.jpg
│       │   └── placeholders/
│       │       ├── avatar.png
│       │       └── room-preview.jpg
│       └── fonts/                 # Custom fonts
│           ├── Inter-Regular.ttf
│           ├── Inter-Medium.ttf
│           ├── Inter-SemiBold.ttf
│           └── Inter-Bold.ttf
├── src/                         # Source code
│   ├── components/               # Reusable components
│   │   ├── common/               # Common components used across the app
│   │   │   ├── Button/           # Button component
│   │   │   │   ├── index.tsx     # Component export
│   │   │   │   ├── Button.tsx    # Button implementation
│   │   │   │   ├── Button.types.ts # Button types
│   │   │   │   └── Button.test.tsx # Button tests
│   │   │   ├── Input/            # Input component
│   │   │   │   ├── index.tsx
│   │   │   │   ├── Input.tsx
│   │   │   │   ├── Input.types.ts
│   │   │   │   └── Input.test.tsx
│   │   │   ├── Modal/            # Modal component
│   │   │   │   ├── index.tsx
│   │   │   │   ├── Modal.tsx
│   │   │   │   ├── Modal.types.ts
│   │   │   │   └── Modal.test.tsx
│   │   │   ├── Loading/          # Loading component
│   │   │   │   ├── index.tsx
│   │   │   │   ├── Loading.tsx
│   │   │   │   └── Loading.test.tsx
│   │   │   ├── Avatar/           # Avatar component
│   │   │   │   ├── index.tsx
│   │   │   │   ├── Avatar.tsx
│   │   │   │   ├── Avatar.types.ts
│   │   │   │   └── Avatar.test.tsx
│   │   │   ├── Badge/            # Badge component
│   │   │   │   ├── index.tsx
│   │   │   │   ├── Badge.tsx
│   │   │   │   ├── Badge.types.ts
│   │   │   │   └── Badge.test.tsx
│   │   │   ├── Tooltip/          # Tooltip component
│   │   │   │   ├── index.tsx
│   │   │   │   ├── Tooltip.tsx
│   │   │   │   ├── Tooltip.types.ts
│   │   │   │   └── Tooltip.test.tsx
│   │   │   └── Dropdown/         # Dropdown component
│   │   │       ├── index.tsx
│   │   │       ├── Dropdown.tsx
│   │   │       ├── Dropdown.types.ts
│   │   │       └── Dropdown.test.tsx
│   │   ├── layout/               # Layout components
│   │   │   ├── Header/           # Header component
│   │   │   │   ├── index.tsx
│   │   │   │   ├── Header.tsx
│   │   │   │   ├── Header.types.ts
│   │   │   │   ├── Header.test.tsx
│   │   │   │   └── Header.css
│   │   │   ├── Sidebar/          # Sidebar component
│   │   │   │   ├── index.tsx
│   │   │   │   ├── Sidebar.tsx
│   │   │   │   ├── Sidebar.types.ts
│   │   │   │   ├── Sidebar.test.tsx
│   │   │   │   └── Sidebar.css
│   │   │   ├── Footer/           # Footer component
│   │   │   │   ├── index.tsx
│   │   │   │   ├── Footer.tsx
│   │   │   │   ├── Footer.types.ts
│   │   │   │   ├── Footer.test.tsx
│   │   │   │   └── Footer.css
│   │   │   └── MainLayout/       # Main layout wrapper
│   │   │       ├── index.tsx
│   │   │       ├── MainLayout.tsx
│   │   │       ├── MainLayout.types.ts
│   │   │       ├── MainLayout.test.tsx
│   │   │       └── MainLayout.css
│   │   ├── forms/                # Form components
│   │   │   ├── LoginForm/        # Login form
│   │   │   │   ├── index.tsx
│   │   │   │   ├── LoginForm.tsx
│   │   │   │   ├── LoginForm.types.ts
│   │   │   │   ├── LoginForm.test.tsx
│   │   │   │   └── LoginForm.css
│   │   │   ├── RegisterForm/     # Register form
│   │   │   │   ├── index.tsx
│   │   │   │   ├── RegisterForm.tsx
│   │   │   │   ├── RegisterForm.types.ts
│   │   │   │   ├── RegisterForm.test.tsx
│   │   │   │   └── RegisterForm.css
│   │   │   ├── CreateRoomForm/   # Create room form
│   │   │   │   ├── index.tsx
│   │   │   │   ├── CreateRoomForm.tsx
│   │   │   │   ├── CreateRoomForm.types.ts
│   │   │   │   ├── CreateRoomForm.test.tsx
│   │   │   │   └── CreateRoomForm.css
│   │   │   ├── ProfileForm/      # Profile form
│   │   │   │   ├── index.tsx
│   │   │   │   ├── ProfileForm.tsx
│   │   │   │   ├── ProfileForm.types.ts
│   │   │   │   ├── ProfileForm.test.tsx
│   │   │   │   └── ProfileForm.css
│   │   │   └── ChangePasswordForm/ # Change password form
│   │   │       ├── index.tsx
│   │   │       ├── ChangePasswordForm.tsx
│   │   │       ├── ChangePasswordForm.types.ts
│   │   │       ├── ChangePasswordForm.test.tsx
│   │   │       └── ChangePasswordForm.css
│   │   └── meeting/              # Meeting-specific components
│   │       ├── VideoPlayer/      # Video player component
│   │       │   ├── index.tsx
│   │       │   ├── VideoPlayer.tsx
│   │       │   ├── VideoPlayer.types.ts
│   │       │   ├── VideoPlayer.test.tsx
│   │       │   └── VideoPlayer.css
│   │       ├── AudioControls/    # Audio controls component
│   │       │   ├── index.tsx
│   │       │   ├── AudioControls.tsx
│   │       │   ├── AudioControls.types.ts
│   │       │   ├── AudioControls.test.tsx
│   │       │   └── AudioControls.css
│   │       ├── ParticipantList/  # Participant list component
│   │       │   ├── index.tsx
│   │       │   ├── ParticipantList.tsx
│   │       │   ├── ParticipantList.types.ts
│   │       │   ├── ParticipantList.test.tsx
│   │       │   └── ParticipantList.css
│   │       ├── ChatBox/          # Chat box component
│   │       │   ├── index.tsx
│   │       │   ├── ChatBox.tsx
│   │       │   ├── ChatBox.types.ts
│   │       │   ├── ChatBox.test.tsx
│   │       │   └── ChatBox.css
│   │       ├── SettingsPanel/    # Settings panel component
│   │       │   ├── index.tsx
│   │       │   ├── SettingsPanel.tsx
│   │       │   ├── SettingsPanel.types.ts
│   │       │   ├── SettingsPanel.test.tsx
│   │       │   └── SettingsPanel.css
│   │       ├── ScreenShare/      # Screen share component
│   │       │   ├── index.tsx
│   │       │   ├── ScreenShare.tsx
│   │       │   ├── ScreenShare.types.ts
│   │       │   ├── ScreenShare.test.tsx
│   │       │   └── ScreenShare.css
│   │       ├── Recording/        # Recording component
│   │       │   ├── index.tsx
│   │       │   ├── Recording.tsx
│   │       │   ├── Recording.types.ts
│   │       │   ├── Recording.test.tsx
│   │       │   └── Recording.css
│   │       └── MeetingControls/  # Meeting controls component
│   │           ├── index.tsx
│   │           ├── MeetingControls.tsx
│   │           ├── MeetingControls.types.ts
│   │           ├── MeetingControls.test.tsx
│   │           └── MeetingControls.css
│   ├── pages/                    # Page components
│   │   ├── Home/                 # Home page
│   │   │   ├── index.tsx         # Home page component
│   │   │   ├── Home.types.ts     # Home page types
│   │   │   ├── Home.test.tsx     # Home page tests
│   │   │   └── Home.css          # Home page styles
│   │   ├── Login/                # Login page
│   │   │   ├── index.tsx
│   │   │   ├── Login.types.ts
│   │   │   ├── Login.test.tsx
│   │   │   └── Login.css
│   │   ├── Register/             # Register page
│   │   │   ├── index.tsx
│   │   │   ├── Register.types.ts
│   │   │   ├── Register.test.tsx
│   │   │   └── Register.css
│   │   ├── Dashboard/            # Dashboard page
│   │   │   ├── index.tsx
│   │   │   ├── Dashboard.types.ts
│   │   │   ├── Dashboard.test.tsx
│   │   │   └── Dashboard.css
│   │   ├── Meeting/              # Meeting page
│   │   │   ├── index.tsx
│   │   │   ├── Meeting.types.ts
│   │   │   ├── Meeting.test.tsx
│   │   │   └── Meeting.css
│   │   ├── Profile/              # Profile page
│   │   │   ├── index.tsx
│   │   │   ├── Profile.types.ts
│   │   │   ├── Profile.test.tsx
│   │   │   └── Profile.css
│   │   ├── Contacts/             # Contacts page
│   │   │   ├── index.tsx
│   │   │   ├── Contacts.types.ts
│   │   │   ├── Contacts.test.tsx
│   │   │   └── Contacts.css
│   │   ├── Rooms/                # Rooms page
│   │   │   ├── index.tsx
│   │   │   ├── Rooms.types.ts
│   │   │   ├── Rooms.test.tsx
│   │   │   └── Rooms.css
│   │   ├── Schedule/             # Schedule page
│   │   │   ├── index.tsx
│   │   │   ├── Schedule.types.ts
│   │   │   ├── Schedule.test.tsx
│   │   │   └── Schedule.css
│   │   ├── Settings/             # Settings page
│   │   │   ├── index.tsx
│   │   │   ├── Settings.types.ts
│   │   │   ├── Settings.test.tsx
│   │   │   └── Settings.css
│   │   ├── Notifications/        # Notifications page
│   │   │   ├── index.tsx
│   │   │   ├── Notifications.types.ts
│   │   │   ├── Notifications.test.tsx
│   │   │   └── Notifications.css
│   │   ├── Recordings/           # Recordings page
│   │   │   ├── index.tsx
│   │   │   ├── Recordings.types.ts
│   │   │   ├── Recordings.test.tsx
│   │   │   └── Recordings.css
│   │   └── NotFound/             # 404 page
│   │       ├── index.tsx
│   │       ├── NotFound.types.ts
│   │       ├── NotFound.test.tsx
│   │       └── NotFound.css
│   ├── stores/                   # State management (Zustand)
│   │   ├── authStore.ts          # Authentication store
│   │   ├── authStore.types.ts    # Authentication store types
│   │   ├── authStore.test.ts    # Authentication store tests
│   │   ├── userStore.ts          # User store
│   │   ├── userStore.types.ts    # User store types
│   │   ├── userStore.test.ts    # User store tests
│   │   ├── roomStore.ts          # Room store
│   │   ├── roomStore.types.ts    # Room store types
│   │   ├── roomStore.test.ts    # Room store tests
│   │   ├── uiStore.ts            # UI store
│   │   ├── uiStore.types.ts      # UI store types
│   │   ├── uiStore.test.ts      # UI store tests
│   │   ├── webrtcStore.ts        # WebRTC store
│   │   ├── webrtcStore.types.ts  # WebRTC store types
│   │   ├── webrtcStore.test.ts  # WebRTC store tests
│   │   ├── notificationStore.ts  # Notification store
│   │   ├── notificationStore.types.ts # Notification store types
│   │   └── notificationStore.test.ts # Notification store tests
│   ├── services/                 # API and service layer
│   │   ├── api/                  # API services
│   │   │   ├── index.ts          # API services export
│   │   │   ├── apiClient.ts      # API client configuration
│   │   │   ├── authApi.ts        # Authentication API
│   │   │   ├── authApi.types.ts  # Authentication API types
│   │   │   ├── authApi.test.ts  # Authentication API tests
│   │   │   ├── userApi.ts        # User API
│   │   │   ├── userApi.types.ts  # User API types
│   │   │   ├── userApi.test.ts  # User API tests
│   │   │   ├── roomApi.ts        # Room API
│   │   │   ├── roomApi.types.ts  # Room API types
│   │   │   ├── roomApi.test.ts  # Room API tests
│   │   │   ├── notificationApi.ts # Notification API
│   │   │   ├── notificationApi.types.ts # Notification API types
│   │   │   ├── notificationApi.test.ts # Notification API tests
│   │   │   ├── webrtcApi.ts      # WebRTC API
│   │   │   ├── webrtcApi.types.ts # WebRTC API types
│   │   │   └── webrtcApi.test.ts # WebRTC API tests
│   │   ├── websocket/            # WebSocket service
│   │   │   ├── index.ts          # WebSocket service export
│   │   │   ├── WebSocketService.ts # WebSocket implementation
│   │   │   ├── WebSocketService.types.ts # WebSocket service types
│   │   │   ├── WebSocketService.test.ts # WebSocket service tests
│   │   │   └── types.ts          # WebSocket types
│   │   ├── webrtc/               # WebRTC service
│   │   │   ├── index.ts          # WebRTC service export
│   │   │   ├── PeerManager.ts    # WebRTC peer manager
│   │   │   ├── PeerManager.types.ts # Peer manager types
│   │   │   ├── PeerManager.test.ts # Peer manager tests
│   │   │   ├── MediaHandler.ts   # Media handler
│   │   │   ├── MediaHandler.types.ts # Media handler types
│   │   │   ├── MediaHandler.test.ts # Media handler tests
│   │   │   └── types.ts          # WebRTC types
│   │   └── storage/              # Storage service
│   │       ├── index.ts          # Storage service export
│   │       ├── LocalStorage.ts   # Local storage implementation
│   │       ├── LocalStorage.types.ts # Local storage types
│   │       ├── LocalStorage.test.ts # Local storage tests
│   │       ├── SessionStorage.ts  # Session storage implementation
│   │       ├── SessionStorage.types.ts # Session storage types
│   │       └── SessionStorage.test.ts # Session storage tests
│   ├── hooks/                    # Custom hooks
│   │   ├── useAuth.ts            # Authentication hook
│   │   ├── useAuth.types.ts      # Authentication hook types
│   │   ├── useAuth.test.ts      # Authentication hook tests
│   │   ├── useUser.ts            # User data hook
│   │   ├── useUser.types.ts      # User data hook types
│   │   ├── useUser.test.ts      # User data hook tests
│   │   ├── useRoom.ts            # Room data hook
│   │   ├── useRoom.types.ts      # Room data hook types
│   │   ├── useRoom.test.ts      # Room data hook tests
│   │   ├── useWebRTC.ts          # WebRTC hook
│   │   ├── useWebRTC.types.ts    # WebRTC hook types
│   │   ├── useWebRTC.test.ts    # WebRTC hook tests
│   │   ├── useWebSocket.ts       # WebSocket hook
│   │   ├── useWebSocket.types.ts # WebSocket hook types
│   │   ├── useWebSocket.test.ts # WebSocket hook tests
│   │   ├── useMedia.ts           # Media devices hook
│   │   ├── useMedia.types.ts     # Media devices hook types
│   │   ├── useMedia.test.ts     # Media devices hook tests
│   │   ├── useLocalStorage.ts    # Local storage hook
│   │   ├── useLocalStorage.types.ts # Local storage hook types
│   │   ├── useLocalStorage.test.ts # Local storage hook tests
│   │   ├── useDebounce.ts        # Debounce hook
│   │   ├── useDebounce.types.ts  # Debounce hook types
│   │   ├── useDebounce.test.ts  # Debounce hook tests
│   │   ├── useThrottle.ts        # Throttle hook
│   │   ├── useThrottle.types.ts  # Throttle hook types
│   │   ├── useThrottle.test.ts  # Throttle hook tests
│   │   ├── useClickOutside.ts    # Click outside hook
│   │   ├── useClickOutside.types.ts # Click outside hook types
│   │   ├── useClickOutside.test.ts # Click outside hook tests
│   │   ├── useKeyPress.ts       # Key press hook
│   │   ├── useKeyPress.types.ts # Key press hook types
│   │   └── useKeyPress.test.ts # Key press hook tests
│   ├── utils/                    # Utility functions
│   │   ├── dateUtils.ts          # Date utilities
│   │   ├── dateUtils.types.ts    # Date utilities types
│   │   ├── dateUtils.test.ts    # Date utilities tests
│   │   ├── validationUtils.ts    # Validation utilities
│   │   ├── validationUtils.types.ts # Validation utilities types
│   │   ├── validationUtils.test.ts # Validation utilities tests
│   │   ├── formatUtils.ts        # Format utilities
│   │   ├── formatUtils.types.ts  # Format utilities types
│   │   ├── formatUtils.test.ts  # Format utilities tests
│   │   ├── constants.ts          # App constants
│   │   ├── config.ts             # App configuration
│   │   ├── config.types.ts       # App configuration types
│   │   ├── config.test.ts       # App configuration tests
│   │   ├── helpers.ts            # Helper functions
│   │   ├── helpers.types.ts      # Helper functions types
│   │   └── helpers.test.ts      # Helper functions tests
│   ├── contexts/                 # React contexts
│   │   ├── ThemeContext.tsx      # Theme context
│   │   ├── ThemeContext.types.ts # Theme context types
│   │   ├── ThemeContext.test.ts # Theme context tests
│   │   ├── AuthContext.tsx       # Auth context
│   │   ├── AuthContext.types.ts # Auth context types
│   │   └── AuthContext.test.ts # Auth context tests
│   ├── styles/                   # Global styles
│   │   ├── global.css            # Global CSS
│   │   ├── tailwind.css          # Tailwind CSS
│   │   ├── theme.css             # Theme-specific styles
│   │   ├── variables.css         # CSS variables
│   │   ├── animations.css        # CSS animations
│   │   └── components.css        # Component styles
│   ├── types/                    # TypeScript type definitions
│   │   ├── auth.ts               # Auth types
│   │   ├── user.ts               # User types
│   │   ├── room.ts               # Room types
│   │   ├── webrtc.ts             # WebRTC types
│   │   ├── notification.ts       # Notification types
│   │   ├── api.ts                # API types
│   │   ├── common.ts             # Common types
│   │   └── index.ts              # Type exports
│   ├── App.tsx                   # Root App component
│   ├── App.types.ts             # App component types
│   ├── App.test.tsx             # App component tests
│   ├── main.tsx                  # App entry point
│   ├── index.css                 # App entry CSS
│   └── vite-env.d.ts            # Vite environment types
├── tests/                       # Test files
│   ├── components/               # Component tests
│   │   ├── common/               # Common component tests
│   │   │   ├── Button.test.tsx
│   │   │   ├── Input.test.tsx
│   │   │   ├── Modal.test.tsx
│   │   │   ├── Loading.test.tsx
│   │   │   ├── Avatar.test.tsx
│   │   │   ├── Badge.test.tsx
│   │   │   ├── Tooltip.test.tsx
│   │   │   └── Dropdown.test.tsx
│   │   ├── layout/               # Layout component tests
│   │   │   ├── Header.test.tsx
│   │   │   ├── Sidebar.test.tsx
│   │   │   ├── Footer.test.tsx
│   │   │   └── MainLayout.test.tsx
│   │   ├── forms/                # Form component tests
│   │   │   ├── LoginForm.test.tsx
│   │   │   ├── RegisterForm.test.tsx
│   │   │   ├── CreateRoomForm.test.tsx
│   │   │   ├── ProfileForm.test.tsx
│   │   │   └── ChangePasswordForm.test.tsx
│   │   └── meeting/              # Meeting component tests
│   │       ├── VideoPlayer.test.tsx
│   │       ├── AudioControls.test.tsx
│   │       ├── ParticipantList.test.tsx
│   │       ├── ChatBox.test.tsx
│   │       ├── SettingsPanel.test.tsx
│   │       ├── ScreenShare.test.tsx
│   │       ├── Recording.test.tsx
│   │       └── MeetingControls.test.tsx
│   ├── pages/                    # Page tests
│   │   ├── Home.test.tsx
│   │   ├── Login.test.tsx
│   │   ├── Register.test.tsx
│   │   ├── Dashboard.test.tsx
│   │   ├── Meeting.test.tsx
│   │   ├── Profile.test.tsx
│   │   ├── Contacts.test.tsx
│   │   ├── Rooms.test.tsx
│   │   ├── Schedule.test.tsx
│   │   ├── Settings.test.tsx
│   │   ├── Notifications.test.tsx
│   │   ├── Recordings.test.tsx
│   │   └── NotFound.test.tsx
│   ├── services/                 # Service tests
│   │   ├── api/                  # API service tests
│   │   │   ├── authApi.test.ts
│   │   │   ├── userApi.test.ts
│   │   │   ├── roomApi.test.ts
│   │   │   ├── notificationApi.test.ts
│   │   │   └── webrtcApi.test.ts
│   │   ├── websocket/            # WebSocket service tests
│   │   │   └── WebSocketService.test.ts
│   │   ├── webrtc/               # WebRTC service tests
│   │   │   ├── PeerManager.test.ts
│   │   │   └── MediaHandler.test.ts
│   │   └── storage/              # Storage service tests
│   │       ├── LocalStorage.test.ts
│   │       └── SessionStorage.test.ts
│   ├── hooks/                    # Hook tests
│   │   ├── useAuth.test.ts
│   │   ├── useUser.test.ts
│   │   ├── useRoom.test.ts
│   │   ├── useWebRTC.test.ts
│   │   ├── useWebSocket.test.ts
│   │   ├── useMedia.test.ts
│   │   ├── useLocalStorage.test.ts
│   │   ├── useDebounce.test.ts
│   │   ├── useThrottle.test.ts
│   │   ├── useClickOutside.test.ts
│   │   └── useKeyPress.test.ts
│   ├── utils/                    # Utility tests
│   │   ├── dateUtils.test.ts
│   │   ├── validationUtils.test.ts
│   │   ├── formatUtils.test.ts
│   │   ├── config.test.ts
│   │   └── helpers.test.ts
│   ├── e2e/                      # End-to-end tests
│   │   ├── auth/                 # Auth e2e tests
│   │   │   ├── login.spec.ts
│   │   │   ├── register.spec.ts
│   │   │   └── logout.spec.ts
│   │   ├── meeting/              # Meeting e2e tests
│   │   │   ├── create-room.spec.ts
│   │   │   ├── join-room.spec.ts
│   │   │   ├── video-controls.spec.ts
│   │   │   ├── chat.spec.ts
│   │   │   └── screen-share.spec.ts
│   │   ├── profile/              # Profile e2e tests
│   │   │   ├── update-profile.spec.ts
│   │   │   ├── change-password.spec.ts
│   │   │   └── contacts.spec.ts
│   │   └── admin/                # Admin e2e tests
│   │       ├── manage-users.spec.ts
│   │       ├── manage-rooms.spec.ts
│   │       └── view-statistics.spec.ts
│   └── setup/                    # Test setup
│       ├── fixtures.ts           # Test fixtures
│       ├── mock-api.ts           # Mock API
│       ├── mock-websocket.ts      # Mock WebSocket
│       └── test-utils.ts         # Test utilities
├── docs/                        # Documentation
│   ├── components/               # Component documentation
│   │   ├── getting-started.md    # Getting started guide
│   │   ├── button.md             # Button documentation
│   │   ├── input.md              # Input documentation
│   │   ├── modal.md              # Modal documentation
│   │   ├── avatar.md             # Avatar documentation
│   │   ├── tooltip.md            # Tooltip documentation
│   │   └── dropdown.md           # Dropdown documentation
│   ├── hooks/                    # Hook documentation
│   │   ├── useAuth.md            # useAuth documentation
│   │   ├── useUser.md            # useUser documentation
│   │   ├── useRoom.md            # useRoom documentation
│   │   ├── useWebRTC.md          # useWebRTC documentation
│   │   ├── useWebSocket.md       # useWebSocket documentation
│   │   ├── useMedia.md           # useMedia documentation
│   │   └── useLocalStorage.md    # useLocalStorage documentation
│   ├── services/                 # Service documentation
│   │   ├── api.md                # API service documentation
│   │   ├── websocket.md          # WebSocket service documentation
│   │   ├── webrtc.md             # WebRTC service documentation
│   │   └── storage.md            # Storage service documentation
│   ├── guides/                   # Developer guides
│   │   ├── theming.md            # Theming guide
│   │   ├── testing.md            # Testing guide
│   │   ├── deployment.md         # Deployment guide
│   │   ├── contributing.md       # Contributing guide
│   │   └── architecture.md       # Architecture overview
│   └── api/                      # API documentation
│       ├── introduction.md       # API introduction
│       ├── authentication.md     # Authentication API
│       ├── users.md              # Users API
│       ├── rooms.md              # Rooms API
│       ├── notifications.md      # Notifications API
│       ├── webrtc.md             # WebRTC API
│       └── websocket.md          # WebSocket API
├── .env.example                 # Environment variables example
├── .eslintrc.js                # ESLint configuration
├── .eslintrc.json              # ESLint configuration
├── .prettierrc.js               # Prettier configuration
├── .prettierrc.json             # Prettier configuration
├── tailwind.config.js            # Tailwind CSS configuration
├── tsconfig.json                # TypeScript configuration
├── tsconfig.node.json           # TypeScript Node configuration
├── vite.config.ts               # Vite configuration
├── package.json                 # Dependencies and scripts
├── postcss.config.js            # PostCSS configuration
├── vitest.config.ts             # Vitest configuration
├── cypress.config.ts            # Cypress configuration
├── components.json              # Component configuration
├── index.html                   # HTML entry point
└── README.md                    # Frontend README
```

## 5. Deployment Structure

```
deployments/
├── docker/                     # Docker configurations
│   ├── Dockerfile.api          # Dockerfile for API server
│   ├── Dockerfile.websocket    # Dockerfile for WebSocket server
│   ├── Dockerfile.media        # Dockerfile for Media server
│   ├── Dockerfile.frontend     # Dockerfile for Frontend
│   ├── Dockerfile.nginx        # Dockerfile for Nginx
│   ├── docker-compose.yml      # Docker Compose for local development
│   └── docker-compose.prod.yml # Docker Compose for production
├── kubernetes/                 # Kubernetes configurations
│   ├── namespaces/             # Namespace configurations
│   │   └── webrtc-meeting.yaml  # Namespace definition
│   ├── configmaps/             # ConfigMap configurations
│   │   ├── api-configmap.yaml   # API server ConfigMap
│   │   ├── websocket-configmap.yaml # WebSocket server ConfigMap
│   │   ├── media-configmap.yaml # Media server ConfigMap
│   │   ├── frontend-configmap.yaml # Frontend ConfigMap
│   │   └── postgres-configmap.yaml # PostgreSQL ConfigMap
│   ├── secrets/                # Secret configurations
│   │   ├── api-secret.yaml      # API server Secret
│   │   ├── websocket-secret.yaml # WebSocket server Secret
│   │   ├── media-secret.yaml    # Media server Secret
│   │   ├── postgres-secret.yaml # PostgreSQL Secret
│   │   └── redis-secret.yaml    # Redis Secret
│   ├── deployments/            # Deployment configurations
│   │   ├── api-deployment.yaml  # API server deployment
│   │   ├── websocket-deployment.yaml # WebSocket server deployment
│   │   ├── media-deployment.yaml # Media server deployment
│   │   ├── frontend-deployment.yaml # Frontend deployment
│   │   ├── postgres-deployment.yaml # PostgreSQL deployment
│   │   └── redis-deployment.yaml # Redis deployment
│   ├── services/               # Service configurations
│   │   ├── api-service.yaml     # API server service
│   │   ├── websocket-service.yaml # WebSocket server service
│   │   ├── media-service.yaml   # Media server service
│   │   ├── frontend-service.yaml # Frontend service
│   │   ├── postgres-service.yaml # PostgreSQL service
│   │   └── redis-service.yaml   # Redis service
│   ├── ingresses/              # Ingress configurations
│   │   ├── api-ingress.yaml     # API server ingress
│   │   ├── websocket-ingress.yaml # WebSocket server ingress
│   │   ├── media-ingress.yaml   # Media server ingress
│   │   └── frontend-ingress.yaml # Frontend ingress
│   ├── persistent-volumes/     # Persistent Volume configurations
│   │   ├── postgres-pv.yaml     # PostgreSQL Persistent Volume
│   │   └── redis-pv.yaml        # Redis Persistent Volume
│   ├── persistent-volume-claims/ # Persistent Volume Claim configurations
│   │   ├── postgres-pvc.yaml    # PostgreSQL Persistent Volume Claim
│   │   └── redis-pvc.yaml       # Redis Persistent Volume Claim
│   ├── horizontal-pod-autoscalers/ # HPA configurations
│   │   ├── api-hpa.yaml         # API server HPA
│   │   ├── websocket-hpa.yaml  # WebSocket server HPA
│   │   ├── media-hpa.yaml      # Media server HPA
│   │   └── frontend-hpa.yaml   # Frontend HPA
│   └── network-policies/       # Network Policy configurations
│       ├── api-network-policy.yaml     # API server network policy
│       ├── websocket-network-policy.yaml # WebSocket server network policy
│       ├── media-network-policy.yaml   # Media server network policy
│       ├── frontend-network-policy.yaml # Frontend network policy
│       ├── postgres-network-policy.yaml # PostgreSQL network policy
│       └── redis-network-policy.yaml   # Redis network policy
├── terraform/                   # Terraform configurations
│   ├── main.tf                 # Main Terraform configuration
│   ├── variables.tf             # Terraform variables
│   ├── outputs.tf              # Terraform outputs
│   ├── provider.tf             # Terraform provider configuration
│   ├── network/                # Network configurations
│   │   ├── vpc.tf              # VPC configuration
│   │   ├── subnet.tf           # Subnet configuration
│   │   ├── security-group.tf   # Security group configuration
│   │   └── internet-gateway.tf # Internet gateway configuration
│   ├── compute/                # Compute configurations
│   │   ├── ec2.tf              # EC2 instances
│   │   ├── ecs.tf              # ECS clusters
│   │   └── eks.tf              # EKS clusters
│   ├── storage/                # Storage configurations
│   │   ├── rds.tf              # RDS instances
│   │   ├── elasticache.tf      # ElastiCache clusters
│   │   └── s3.tf               # S3 buckets
│   └── cdn/                    # CDN configurations
│       └── cloudfront.tf       # CloudFront distribution
└── ansible/                    # Ansible configurations
    ├── inventory.ini           # Ansible inventory
    ├── playbook.yml            # Main playbook
    ├── roles/                  # Ansible roles
    │   ├── common/             # Common role
    │   ├── docker/             # Docker role
    │   ├── nginx/              # Nginx role
    │   ├── postgres/           # PostgreSQL role
    │   ├── redis/              # Redis role
    │   ├── api/                # API server role
    │   ├── websocket/          # WebSocket server role
    │   ├── media/              # Media server role
    │   └── frontend/           # Frontend role
    └── templates/              # Ansible templates
        ├── docker-compose.yml.j2
        ├── nginx.conf.j2
        ├── config.yaml.j2
        └── env.j2
```

## 6. Documentation Structure

```
docs/
├── README.md                   # Documentation index
├── getting-started/           # Getting started guides
│   ├── installation.md         # Installation guide
│   ├── setup.md               # Setup guide
│   ├── configuration.md       # Configuration guide
│   └── deployment.md         # Deployment guide
├── architecture/              # Architecture documentation
│   ├── overview.md            # Architecture overview
│   ├── backend.md             # Backend architecture
│   ├── frontend.md            # Frontend architecture
│   ├── database.md            # Database architecture
│   ├── api.md                 # API architecture
│   ├── webrtc.md              # WebRTC architecture
│   └── deployment.md          # Deployment architecture
├── development/              # Development guides
│   ├── environment-setup.md   # Development environment setup
│   ├── coding-standards.md    # Coding standards
│   ├── git-workflow.md        # Git workflow
│   ├── testing.md             # Testing guide
│   ├── debugging.md           # Debugging guide
│   └── contributing.md       # Contributing guide
├── api/                       # API documentation
│   ├── introduction.md         # API introduction
│   ├── authentication.md       # Authentication API
│   ├── users.md                # Users API
│   ├── rooms.md                # Rooms API
│   ├── notifications.md       # Notifications API
│   ├── webrtc.md               # WebRTC API
│   ├── websocket.md           # WebSocket API
│   └── errors.md               # Error handling
├── components/                # Component documentation
│   ├── introduction.md         # Component introduction
│   ├── common/                 # Common components
│   │   ├── button.md           # Button component
│   │   ├── input.md            # Input component
│   │   ├── modal.md            # Modal component
│   │   ├── avatar.md           # Avatar component
│   │   ├── badge.md            # Badge component
│   │   ├── tooltip.md          # Tooltip component
│   │   └── dropdown.md         # Dropdown component
│   ├── layout/                 # Layout components
│   │   ├── header.md           # Header component
│   │   ├── sidebar.md          # Sidebar component
│   │   ├── footer.md           # Footer component
│   │   └── main-layout.md       # Main layout component
│   ├── forms/                  # Form components
│   │   ├── login-form.md       # Login form component
│   │   ├── register-form.md    # Register form component
│   │   ├── create-room-form.md # Create room form component
│   │   ├── profile-form.md     # Profile form component
│   │   └── change-password-form.md # Change password form component
│   └── meeting/                # Meeting components
│       ├── video-player.md     # Video player component
│       ├── audio-controls.md   # Audio controls component
│       ├── participant-list.md # Participant list component
│       ├── chat-box.md         # Chat box component
│       ├── settings-panel.md   # Settings panel component
│       ├── screen-share.md     # Screen share component
│       ├── recording.md        # Recording component
│       └── meeting-controls.md # Meeting controls component
├── hooks/                     # Hook documentation
│   ├── introduction.md         # Hook introduction
│   ├── useAuth.md              # useAuth hook
│   ├── useUser.md              # useUser hook
│   ├── useRoom.md              # useRoom hook
│   ├── useWebRTC.md            # useWebRTC hook
│   ├── useWebSocket.md         # useWebSocket hook
│   ├── useMedia.md             # useMedia hook
│   ├── useLocalStorage.md     # useLocalStorage hook
│   ├── useDebounce.md          # useDebounce hook
│   ├── useThrottle.md          # useThrottle hook
│   ├── useClickOutside.md      # useClickOutside hook
│   └── useKeyPress.md          # useKeyPress hook
├── services/                  # Service documentation
│   ├── introduction.md         # Service introduction
│   ├── api/                    # API services
│   │   ├── auth-api.md         # Authentication API service
│   │   ├── user-api.md         # User API service
│   │   ├── room-api.md         # Room API service
│   │   ├── notification-api.md # Notification API service
│   │   └── webrtc-api.md       # WebRTC API service
│   ├── websocket/              # WebSocket service
│   │   ├── webSocket-service.md # WebSocket service
│   │   ├── types.md             # WebSocket types
│   │   └── events.md           # WebSocket events
│   ├── webrtc/                 # WebRTC service
│   │   ├── peer-manager.md      # Peer manager
│   │   ├── media-handler.md    # Media handler
│   │   ├── types.md             # WebRTC types
│   │   └── signaling.md         # WebRTC signaling
│   └── storage/                # Storage service
│       ├── local-storage.md    # Local storage
│       ├── session-storage.md  # Session storage
│       └── types.md             # Storage types
├── deployment/                # Deployment documentation
│   ├── overview.md            # Deployment overview
│   ├── docker.md               # Docker deployment
│   ├── kubernetes.md          # Kubernetes deployment
│   ├── terraform.md            # Terraform deployment
│   ├── ansible.md              # Ansible deployment
│   ├── monitoring.md           # Monitoring setup
│   └── scaling.md             # Scaling guide
├── testing/                   # Testing documentation
│   ├── overview.md            # Testing overview
│   ├── unit-testing.md         # Unit testing
│   ├── integration-testing.md  # Integration testing
│   ├── e2e-testing.md         # End-to-end testing
│   ├── api-testing.md         # API testing
│   ├── performance-testing.md  # Performance testing
│   └── security-testing.md    # Security testing
├── security/                  # Security documentation
│   ├── overview.md            # Security overview
│   ├── authentication.md       # Authentication security
│   ├── authorization.md        # Authorization security
│   ├── data-protection.md     # Data protection
│   ├── network-security.md    # Network security
│   ├── api-security.md        # API security
│   └── webRTC-security.md     # WebRTC security
├── monitoring/                # Monitoring documentation
│   ├── overview.md            # Monitoring overview
│   ├── logging.md             # Logging setup
│   ├── metrics.md             # Metrics collection
│   ├── tracing.md             # Distributed tracing
│   ├── alerting.md            # Alerting setup
│   └── dashboard.md           # Dashboard setup
├── troubleshooting/            # Troubleshooting guides
│   ├── common-issues.md       # Common issues
│   ├── performance-issues.md  # Performance issues
│   ├── connection-issues.md   # Connection issues
│   ├── deployment-issues.md   # Deployment issues
│   └── debugging-tips.md      # Debugging tips
├── release-notes/             # Release notes
│   ├── v1.0.0.md              # Version 1.0.0 release notes
│   ├── v1.1.0.md              # Version 1.1.0 release notes
│   └── v2.0.0.md              # Version 2.0.0 release notes
└── glossary/                  # Glossary of terms
    ├── api-terms.md           # API terms
    ├── frontend-terms.md      # Frontend terms
    ├── backend-terms.md       # Backend terms
    ├── database-terms.md      # Database terms
    ├── deployment-terms.md    # Deployment terms
    └── webrtc-terms.md        # WebRTC terms
```

## 7. Configuration Files

### 7.1 Root Configuration Files

```
# .env.example
# Backend Configuration
BACKEND_PORT=8080
BACKEND_HOST=localhost
DATABASE_URL=postgresql://user:password@localhost:5432/webrtc_meeting?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-jwt-secret-key
JWT_REFRESH_SECRET=your-jwt-refresh-secret-key

# Frontend Configuration
FRONTEND_PORT=3000
FRONTEND_HOST=localhost
API_BASE_URL=http://localhost:8080/api
WEBSOCKET_URL=ws://localhost:8080/ws

# WebRTC Configuration
STUN_SERVERS=stun:stun.l.google.com:19302,stun:stun1.l.google.com:19302
TURN_SERVERS=turn:your-turn-server.com:3478
TURN_USERNAME=your-turn-username
TURN_CREDENTIAL=your-turn-credential

# Email Configuration
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your-smtp-user
SMTP_PASSWORD=your-smtp-password
SMTP_FROM=noreply@example.com

# Monitoring Configuration
LOG_LEVEL=info
METRICS_ENABLED=true
TRACING_ENABLED=true

# Storage Configuration
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=10485760  # 10MB
ALLOWED_FILE_TYPES=jpg,jpeg,png,gif,pdf,doc,docx
```

```
# .gitignore
# Dependencies
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Production builds
dist/
build/

# Environment variables
.env
.env.local
.env.development.local
.env.test.local
.env.production.local

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Logs
logs/
*.log

# Database
*.sqlite
*.db

# Coverage
coverage/
.nyc_output/

# Temporary files
tmp/
temp/

# Docker
.dockerignore

# Go
go.sum
backend/bin/
backend/pkg/

# Test
backend/.test/
backend/coverage.out
frontend/coverage/
frontend/.nyc_output/

# Preact
frontend/dist/
frontend/.vite/

# Cypress
frontend/cypress/videos/
frontend/cypress/screenshots/

# Storybook
frontend/storybook-static/
```

```
# docker-compose.yml
version: '3.8'

services:
  # PostgreSQL Database
  postgres:
    image: postgres:13
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: webrtc_meeting
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/database/migrations:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"
    networks:
      - webrtc-network

  # Redis Cache
  redis:
    image: redis:6-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - webrtc-network

  # API Server
  api:
    build:
      context: ./backend
      dockerfile: ./deployments/docker/Dockerfile.api
    environment:
      - DATABASE_URL=postgresql://user:password@postgres:5432/webrtc_meeting?sslmode=disable
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=your-jwt-secret-key
      - JWT_REFRESH_SECRET=your-jwt-refresh-secret-key
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    networks:
      - webrtc-network
    volumes:
      - ./backend:/app
      - /app/node_modules

  # WebSocket Server
  websocket:
    build:
      context: ./backend
      dockerfile: ./deployments/docker/Dockerfile.websocket
    environment:
      - DATABASE_URL=postgresql://user:password@postgres:5432/webrtc_meeting?sslmode=disable
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=your-jwt-secret-key
    ports:
      - "8081:8081"
    depends_on:
      - postgres
      - redis
    networks:
      - webrtc-network
    volumes:
      - ./backend:/app
      - /app/node_modules

  # Media Server
  media:
    build:
      context: ./backend
      dockerfile: ./deployments/docker/Dockerfile.media
    environment:
      - DATABASE_URL=postgresql://user:password@postgres:5432/webrtc_meeting?sslmode=disable
      - REDIS_URL=redis://redis:6379
      - STUN_SERVERS=stun:stun.l.google.com:19302
      - TURN_SERVERS=turn:your-turn-server.com:3478
      - TURN_USERNAME=your-turn-username
      - TURN_CREDENTIAL=your-turn-credential
    ports:
      - "8082:8082"
      - "10000-10100:10000-10100/udp" # UDP ports for WebRTC
    depends_on:
      - postgres
      - redis
    networks:
      - webrtc-network
    volumes:
      - ./backend:/app
      - /app/node_modules

  # Frontend
  frontend:
    build:
      context: ./frontend
      dockerfile: ./deployments/docker/Dockerfile.frontend
    environment:
      - API_BASE_URL=http://localhost:8080/api
      - WEBSOCKET_URL=ws://localhost:8081/ws
    ports:
      - "3000:3000"
    depends_on:
      - api
      - websocket
    networks:
      - webrtc-network
    volumes:
      - ./frontend:/app
      - /app/node_modules

  # Nginx Reverse Proxy
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./deployments/docker/nginx.conf:/etc/nginx/nginx.conf
      - ./deployments/docker/ssl:/etc/nginx/ssl
    depends_on:
      - api
      - websocket
      - media
      - frontend
    networks:
      - webrtc-network

volumes:
  postgres_data:
  redis_data:

networks:
  webrtc-network:
    driver: bridge
```

### 7.2 Backend Configuration Files

```
# backend/go.mod
module github.com/your-org/webrtc-meeting-backend

go 1.21

require (
  github.com/gin-gonic/gin v1.9.1
  github.com/gorilla/websocket v1.5.0
  github.com/pion/webrtc/v3 v3.2.24
  github.com/joho/godotenv v1.4.0
  gorm.io/gorm v1.25.2
  gorm.io/driver/postgres v1.5.2
  github.com/golang-jwt/jwt/v5 v5.0.0
  github.com/go-redis/redis/v8 v8.11.5
  github.com/sirupsen/logrus v1.9.3
  github.com/prometheus/client_golang v1.16.0
  github.com/opentracing/opentracing-go v1.2.0
  github.com/uber/jaeger-client-go v2.30.0+incompatible
  github.com/testcontainers/testcontainers-go v0.20.0
  github.com/stretchr/testify v1.8.4
  github.com/golang/mock v1.6.0
  github.com/pressly/goose/v3 v3.12.0
  github.com/swaggo/swag v1.16.2
  github.com/swaggo/gin-swagger v1.6.0
  github.com/go-playground/validator/v10 v10.14.0
  github.com/casbin/casbin/v2 v2.82.0
  github.com/ulule/deepcopier v0.0.0-20200417032501-5fac5c633592
  github.com/google/uuid v1.3.0
  golang.org/x/crypto v0.12.0
  golang.org/x/time v0.3.0
)

require (
  github.com/bytedance/sonic v1.9.1 // indirect
  github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
  github.com/gabriel-vasile/mimetype v1.4.2 // indirect
  github.com/gin-contrib/sse v0.1.0 // indirect
  github.com/go-playground/locales v0.14.1 // indirect
  github.com/go-playground/universal-translator v0.18.1 // indirect
  github.com/goccy/go-json v0.10.2 // indirect
  github.com/jackc/pgpassfile v1.0.0 // indirect
  github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
  github.com/jackc/pgx/v5 v5.4.3 // indirect
  github.com/jinzhu/inflection v1.0.0 // indirect
  github.com/jinzhu/now v1.1.5 // indirect
  github.com/json-iterator/go v1.1.12 // indirect
  github.com/klauspost/cpuid/v2 v2.2.4 // indirect
  github.com/leodido/go-urn v1.2.4 // indirect
  github.com/mattn/go-isatty v0.0.19 // indirect
  github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
  github.com/modern-go/reflect2 v1.0.2 // indirect
  github.com/pelletier/go-toml/v2 v2.0.8 // indirect
  github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
  github.com/ugorji/go/codec v1.2.11 // indirect
  golang.org/x/arch v0.3.0 // indirect
  golang.org/x/net v0.10.0 // indirect
  golang.org/x/sys v0.11.0 // indirect
  golang.org/x/text v0.12.0 // indirect
  google.golang.org/protobuf v1.30.0 // indirect
  gopkg.in/yaml.v3 v3.0.1 // indirect
)
```

```
# backend/Makefile
.PHONY: build run test clean migrate-up migrate-down migrate-force deps lint fmt vet

# Variables
BINARY_NAME=webrtc-meeting-api
BINARY_UNIX=$(BINARY_NAME)_unix
VERSION=1.0.0
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -w -s"

# Default target
all: clean deps lint test build

# Build the application
build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/api

# Build for Unix
build-unix:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_UNIX) ./cmd/api

# Run the application
run:
	go run ./cmd/api

# Run tests
test:
	go test -v -race -cover=./... ./...

# Run tests with coverage
test-coverage:
	go test -v -race -cover=./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run integration tests
test-integration:
	go test -v -tags=integration ./tests/integration/...

# Run e2e tests
test-e2e:
	go test -v -tags=e2e ./tests/e2e/...

# Clean build artifacts
clean:
	go clean
	rm -rf bin/

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run vet
vet:
	go vet ./...

# Run migrations up
migrate-up:
	goose -dir=./backend/database/migrations postgres "user=user password=password dbname=webrtc_meeting sslmode=disable" up

# Run migrations down
migrate-down:
	goose -dir=./backend/database/migrations postgres "user=user password=password dbname=webrtc_meeting sslmode=disable" down

# Force run migrations
migrate-force:
	goose -dir=./backend/database/migrations postgres "user=user password=password dbname=webrtc_meeting sslmode=disable" force

# Create new migration
create-migration:
	@read -p "Enter migration name: " name; \
	goose -dir=./backend/database/migrations create $$name sql

# Build Docker image
docker-build:
	docker build -f ./deployments/docker/Dockerfile.api -t webrtc-meeting-api:$(VERSION) .

# Push Docker image
docker-push:
	docker push webrtc-meeting-api:$(VERSION)

# Run Docker Compose
docker-up:
	docker-compose -f ./deployments/docker/docker-compose.yml up -d

# Stop Docker Compose
docker-down:
	docker-compose -f ./deployments/docker/docker-compose.yml down

# View Docker logs
docker-logs:
	docker-compose -f ./deployments/docker/docker-compose.yml logs -f

# Generate API documentation
docs:
	swag init -g ./cmd/api/main.go -o ./docs/api

# Run development server with hot reload
dev:
	air -c .air.toml

# Install development tools
install-dev-tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install -tags 'postgres' github.com/pressly/goose/v3/cmd/goose@latest
```

### 7.3 Frontend Configuration Files

```
# frontend/package.json
{
  "name": "webrtc-meeting-frontend",
  "private": true,
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview",
    "test": "vitest",
    "test:ui": "vitest --ui",
    "test:coverage": "vitest --coverage",
    "test:e2e": "cypress run",
    "test:e2e:open": "cypress open",
    "lint": "eslint . --ext ts,tsx --report-unused-disable-directives --max-warnings 0",
    "lint:fix": "eslint . --ext ts,tsx --report-unused-disable-directives --max-warnings 0 --fix",
    "format": "prettier --write \"src/**/*.{ts,tsx,css,md}\"",
    "format:check": "prettier --check \"src/**/*.{ts,tsx,css,md}\"",
    "typecheck": "tsc --noEmit",
    "storybook": "storybook dev -p 6006",
    "build-storybook": "storybook build"
  },
  "dependencies": {
    "preact": "^10.17.1",
    "preact-iso": "^1.1.0",
    "preact-render-to-string": "^6.3.1",
    "@preact/preset-vite": "^2.5.0",
    "preact-router": "^4.1.0",
    "zustand": "^4.4.1",
    "axios": "^1.5.0",
    "simple-peer": "^9.11.1",
    "socket.io-client": "^4.7.2",
    "tailwindcss": "^3.3.3",
    "clsx": "^2.0.0",
    "tailwind-merge": "^1.14.0",
    "class-variance-authority": "^0.7.0",
    "lucide-preact": "^0.288.0",
    "@radix-ui/react-slot": "^1.0.2",
    "@radix-ui/react-dropdown-menu": "^2.0.6",
    "@radix-ui/react-dialog": "^1.0.5",
    "@radix-ui/react-toast": "^1.1.5",
    "react-hook-form": "^7.46.2",
    "@hookform/resolvers": "^3.3.1",
    "zod": "^3.22.2",
    "date-fns": "^2.30.0",
    "framer-motion": "^10.16.4"
  },
  "devDependencies": {
    "@types/node": "^20.5.7",
    "@types/react": "^18.2.15",
    "@types/react-dom": "^18.2.7",
    "@typescript-eslint/eslint-plugin": "^6.0.0",
    "@typescript-eslint/parser": "^6.0.0",
    "eslint": "^8.45.0",
    "eslint-plugin-react-hooks": "^4.6.0",
    "eslint-plugin-react-refresh": "^0.4.3",
    "typescript": "^5.0.2",
    "vite": "^4.4.5",
    "vitest": "^0.34.4",
    "@vitest/ui": "^0.34.4",
    "@vitest/coverage-v8": "^0.34.4",
    "jsdom": "^22.1.0",
    "cypress": "^13.3.0",
    "@testing-library/preact": "^3.2.3",
    "@testing-library/jest-dom": "^6.1.3",
    "@testing-library/user-event": "^14.4.3",
    "prettier": "^3.0.2",
    "eslint-config-prettier": "^9.0.0",
    "eslint-plugin-prettier": "^5.0.0",
    "autoprefixer": "^10.4.15",
    "postcss": "^8.4.29",
    "storybook": "^7.4.5",
    "@storybook/addon-essentials": "^7.4.5",
    "@storybook/addon-interactions": "^7.4.5",
    "@storybook/addon-links": "^7.4.5",
    "@storybook/blocks": "^7.4.5",
    "@storybook/preact": "^7.4.5",
    "@storybook/preact-vite": "^7.4.5",
    "@storybook/testing-library": "^0.2.0"
  }
}
```

```
# frontend/vite.config.ts
import { defineConfig } from 'vite'
import preact from '@preact/preset-vite'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [preact()],
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
      '@components': resolve(__dirname, './src/components'),
      '@pages': resolve(__dirname, './src/pages'),
      '@stores': resolve(__dirname, './src/stores'),
      '@services': resolve(__dirname, './src/services'),
      '@hooks': resolve(__dirname, './src/hooks'),
      '@utils': resolve(__dirname, './src/utils'),
      '@types': resolve(__dirname, './src/types'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/ws': {
        target: 'ws://localhost:8081',
        ws: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['preact', 'preact-router'],
          state: ['zustand'],
          ui: ['tailwindcss', 'clsx', 'tailwind-merge'],
          forms: ['react-hook-form', '@hookform/resolvers', 'zod'],
          icons: ['lucide-preact'],
          utils: ['date-fns', 'clsx'],
        },
      },
    },
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./tests/setup/test-utils.ts'],
  },
})
```

```
# frontend/tailwind.config.js
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Primary Colors (Hijau)
        primary: {
          50: '#f0fdf4',
          100: '#dcfce7',
          200: '#bbf7d0',
          300: '#86efac',
          400: '#4ade80',
          500: '#22c55e', // Primary Green
          600: '#16a34a',
          700: '#15803d',
          800: '#166534',
          900: '#14532d',
        },
        // Secondary Colors (Abu-abu)
        secondary: {
          50: '#f9fafb',
          100: '#f3f4f6',
          200: '#e5e7eb',
          300: '#d1d5db',
          400: '#9ca3af',
          500: '#6b7280',
          600: '#4b5563',
          700: '#374151',
          800: '#1f2937',
          900: '#111827',
        },
        // Accent Colors
        accent: {
          500: '#10b981',
          600: '#059669',
        },
        // Status colors
        success: '#22c55e',
        warning: '#f59e0b',
        error: '#ef4444',
        info: '#3b82f6',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      boxShadow: {
        'custom': '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
        'custom-lg': '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
      },
      animation: {
        'fade-in': 'fadeIn 0.5s ease-in-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
    },
  },
  plugins: [],
}
```

## 8. Kesimpulan

Struktur folder proyek untuk aplikasi WebRTC meeting ini dirancang dengan mempertimbangkan:

1. **Organisasi yang jelas** antara backend (Golang) dan frontend (Preact)
2. **Separation of concerns** dengan pemisahan yang jelas antara components, services, dan utilities
3. **Scalability** dengan struktur yang dapat menangani pertumbuhan fitur dan tim
4. **Maintainability** dengan konvensi penamaan dan organisasi file yang konsisten
5. **Testability** dengan struktur yang memudahkan penulisan tes di berbagai level
6. **Deployment readiness** dengan konfigurasi Docker, Kubernetes, dan dokumentasi yang lengkap
7. **Developer experience** dengan tooling dan automation yang memudahkan pengembangan

Dengan struktur folder ini, tim pengembang dapat dengan mudah menavigasi kode, memahami arsitektur, dan berkolaborasi secara efektif dalam membangun aplikasi WebRTC meeting.