# Desain Database/Schema untuk User Management dan Room System

## 1. Pendahuluan

Desain database untuk aplikasi WebRTC meeting ini dirancang untuk mendukung fitur-fitur utama seperti user management, room-based meeting system, dan riwayat meeting. Database akan menggunakan PostgreSQL sebagai RDBMS utama dengan Redis untuk caching dan session management.

## 2. Pilihan Database dan Alasan

### 2.1 PostgreSQL sebagai Database Utama

**Alasan memilih PostgreSQL:**
- **Open-source dan mature**: PostgreSQL adalah database open-source yang telah matang dan stabil
- **Fitur lengkap**: Mendukung fitur-fitur seperti JSONB, full-text search, dan advanced indexing
- **Scalability**: Dapat menangani beban yang tinggi dengan baik
- **Reliability**: ACID compliance dan data integrity yang kuat
- **Ekosistem**: Banyak tool dan library yang mendukung PostgreSQL
- **Performance**: Performa yang baik untuk query kompleks dan concurrent access

### 2.2 Redis untuk Caching dan Session Management

**Alasan memilih Redis:**
- **In-memory storage**: Performa sangat tinggi untuk data yang sering diakses
- **Data structures**: Mendukung berbagai tipe data structures (strings, hashes, lists, sets, etc.)
- **Pub/Sub**: Mendukung messaging pattern untuk real-time updates
- **Persistence**: Opsi untuk menyimpan data ke disk jika diperlukan
- **Scalability**: Dapat di-scale dengan Redis Cluster

## 3. Entity Relationship Diagram (ERD)

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│     users       │       │      rooms      │       │  room_participants│
├─────────────────┤       ├─────────────────┤       ├─────────────────┤
│ id (PK)         │───┐   │ id (PK)         │───┐   │ id (PK)         │
│ email           │   │   │ name            │   │   │ room_id (FK)    │
│ password        │   │   │ description     │   │   │ user_id (FK)    │
│ first_name      │   │   │ host_id (FK)    │◄──┘   │ joined_at       │
│ last_name       │   │   │ password        │       │ left_at         │
│ avatar          │   │   │ max_users       │       │ role            │
│ status          │   │   │ status          │       └─────────────────┘
│ created_at      │   │   │ created_at      │
│ updated_at      │   │   │ updated_at      │
└─────────────────┘   │   │ ended_at        │
                      │   └─────────────────┘
                      │
                      │   ┌─────────────────┐
                      │   │  user_contacts  │
                      │   ├─────────────────┤
                      └───│ user_id (FK)    │
                          │ contact_id (FK) │
                          │ created_at      │
                          └─────────────────┘

┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│   user_sessions │       │  room_messages  │       │ meeting_history │
├─────────────────┤       ├─────────────────┤       ├─────────────────┤
│ id (PK)         │       │ id (PK)         │       │ id (PK)         │
│ user_id (FK)    │       │ room_id (FK)    │       │ room_id (FK)    │
│ token           │       │ user_id (FK)    │       │ user_id (FK)    │
│ expires_at      │       │ message         │       │ joined_at       │
│ created_at      │       │ message_type    │       │ left_at         │
│ ip_address      │       │ created_at      │       │ duration        │
│ user_agent      │       └─────────────────┘       └─────────────────┘
└─────────────────┘

┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│  user_settings  │       │  room_settings  │       │  notifications  │
├─────────────────┤       ├─────────────────┤       ├─────────────────┤
│ id (PK)         │       │ id (PK)         │       │ id (PK)         │
│ user_id (FK)    │       │ room_id (FK)    │       │ user_id (FK)    │
│ setting_key     │       │ setting_key     │       │ type            │
│ setting_value   │       │ setting_value   │       │ title           │
│ created_at      │       │ created_at      │       │ message         │
│ updated_at      │       │ updated_at      │       │ is_read         │
└─────────────────┘       └─────────────────┘       │ created_at      │
                                                    └─────────────────┘
```

## 4. Tabel Database Detail

### 4.1 Tabel Users

Tabel `users` menyimpan data pengguna aplikasi.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    avatar TEXT,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'banned')),
    email_verified BOOLEAN DEFAULT FALSE,
    last_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Comments
COMMENT ON TABLE users IS 'Table to store user information';
COMMENT ON COLUMN users.id IS 'Unique identifier for the user';
COMMENT ON COLUMN users.email IS 'User email address (unique)';
COMMENT ON COLUMN users.password IS 'Hashed password';
COMMENT ON COLUMN users.first_name IS 'User first name';
COMMENT ON COLUMN users.last_name IS 'User last name';
COMMENT ON COLUMN users.avatar IS 'URL to user avatar image';
COMMENT ON COLUMN users.status IS 'User status (active, inactive, banned)';
COMMENT ON COLUMN users.email_verified IS 'Flag indicating if email is verified';
COMMENT ON COLUMN users.last_login IS 'Timestamp of last login';
COMMENT ON COLUMN users.created_at IS 'Timestamp when user was created';
COMMENT ON COLUMN users.updated_at IS 'Timestamp when user was last updated';
```

### 4.2 Tabel Rooms

Tabel `rooms` menyimpan data ruang meeting.

```sql
CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    host_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    password VARCHAR(255),
    max_users INTEGER DEFAULT 10 CHECK (max_users > 0),
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'ended', 'cancelled')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ended_at TIMESTAMP WITH TIME ZONE,
    scheduled_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_rooms_host_id ON rooms(host_id);
CREATE INDEX idx_rooms_status ON rooms(status);
CREATE INDEX idx_rooms_created_at ON rooms(created_at);
CREATE INDEX idx_rooms_scheduled_at ON rooms(scheduled_at);

-- Comments
COMMENT ON TABLE rooms IS 'Table to store meeting rooms';
COMMENT ON COLUMN rooms.id IS 'Unique identifier for the room';
COMMENT ON COLUMN rooms.name IS 'Room name';
COMMENT ON COLUMN rooms.description IS 'Room description';
COMMENT ON COLUMN rooms.host_id IS 'ID of the room host (references users.id)';
COMMENT ON COLUMN rooms.password IS 'Room password (optional)';
COMMENT ON COLUMN rooms.max_users IS 'Maximum number of users allowed in the room';
COMMENT ON COLUMN rooms.status IS 'Room status (active, ended, cancelled)';
COMMENT ON COLUMN rooms.created_at IS 'Timestamp when room was created';
COMMENT ON COLUMN rooms.updated_at IS 'Timestamp when room was last updated';
COMMENT ON COLUMN rooms.ended_at IS 'Timestamp when room was ended';
COMMENT ON COLUMN rooms.scheduled_at IS 'Timestamp when room is scheduled to start';
```

### 4.3 Tabel Room Participants

Tabel `room_participants` menyimpan data peserta dalam ruang meeting.

```sql
CREATE TABLE room_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    left_at TIMESTAMP WITH TIME ZONE,
    role VARCHAR(50) DEFAULT 'participant' CHECK (role IN ('host', 'participant', 'moderator')),
    
    -- Unique constraint to prevent duplicate participants
    UNIQUE(room_id, user_id)
);

-- Indexes
CREATE INDEX idx_room_participants_room_id ON room_participants(room_id);
CREATE INDEX idx_room_participants_user_id ON room_participants(user_id);
CREATE INDEX idx_room_participants_joined_at ON room_participants(joined_at);

-- Comments
COMMENT ON TABLE room_participants IS 'Table to store room participants';
COMMENT ON COLUMN room_participants.id IS 'Unique identifier for the participant';
COMMENT ON COLUMN room_participants.room_id IS 'ID of the room (references rooms.id)';
COMMENT ON COLUMN room_participants.user_id IS 'ID of the user (references users.id)';
COMMENT ON COLUMN room_participants.joined_at IS 'Timestamp when user joined the room';
COMMENT ON COLUMN room_participants.left_at IS 'Timestamp when user left the room';
COMMENT ON COLUMN room_participants.role IS 'Role of the user in the room (host, participant, moderator)';
```

### 4.4 Tabel User Contacts

Tabel `user_contacts` menyimpan data kontak pengguna.

```sql
CREATE TABLE user_contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contact_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Unique constraint to prevent duplicate contacts
    UNIQUE(user_id, contact_id),
    
    -- Check constraint to prevent user from adding themselves as contact
    CHECK(user_id != contact_id)
);

-- Indexes
CREATE INDEX idx_user_contacts_user_id ON user_contacts(user_id);
CREATE INDEX idx_user_contacts_contact_id ON user_contacts(contact_id);

-- Comments
COMMENT ON TABLE user_contacts IS 'Table to store user contacts';
COMMENT ON COLUMN user_contacts.id IS 'Unique identifier for the contact';
COMMENT ON COLUMN user_contacts.user_id IS 'ID of the user (references users.id)';
COMMENT ON COLUMN user_contacts.contact_id IS 'ID of the contact (references users.id)';
COMMENT ON COLUMN user_contacts.created_at IS 'Timestamp when contact was added';
```

### 4.5 Tabel User Sessions

Tabel `user_sessions` menyimpan data sesi pengguna untuk autentikasi.

```sql
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT,
    is_revoked BOOLEAN DEFAULT FALSE
);

-- Indexes
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_token ON user_sessions(token);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

-- Comments
COMMENT ON TABLE user_sessions IS 'Table to store user sessions for authentication';
COMMENT ON COLUMN user_sessions.id IS 'Unique identifier for the session';
COMMENT ON COLUMN user_sessions.user_id IS 'ID of the user (references users.id)';
COMMENT ON COLUMN user_sessions.token IS 'Session token (JWT)';
COMMENT ON COLUMN user_sessions.expires_at IS 'Timestamp when session expires';
COMMENT ON COLUMN user_sessions.created_at IS 'Timestamp when session was created';
COMMENT ON COLUMN user_sessions.ip_address IS 'IP address of the user';
COMMENT ON COLUMN user_sessions.user_agent IS 'User agent string';
COMMENT ON COLUMN user_sessions.is_revoked IS 'Flag indicating if session is revoked';
```

### 4.6 Tabel Room Messages

Tabel `room_messages` menyimpan pesan chat dalam ruang meeting.

```sql
CREATE TABLE room_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    message_type VARCHAR(50) DEFAULT 'text' CHECK (message_type IN ('text', 'system', 'file')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_room_messages_room_id ON room_messages(room_id);
CREATE INDEX idx_room_messages_user_id ON room_messages(user_id);
CREATE INDEX idx_room_messages_created_at ON room_messages(created_at);

-- Comments
COMMENT ON TABLE room_messages IS 'Table to store room chat messages';
COMMENT ON COLUMN room_messages.id IS 'Unique identifier for the message';
COMMENT ON COLUMN room_messages.room_id IS 'ID of the room (references rooms.id)';
COMMENT ON COLUMN room_messages.user_id IS 'ID of the user who sent the message (references users.id)';
COMMENT ON COLUMN room_messages.message IS 'Message content';
COMMENT ON COLUMN room_messages.message_type IS 'Type of message (text, system, file)';
COMMENT ON COLUMN room_messages.created_at IS 'Timestamp when message was sent';
```

### 4.7 Tabel Meeting History

Tabel `meeting_history` menyimpan riwayat meeting yang telah selesai.

```sql
CREATE TABLE meeting_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL,
    left_at TIMESTAMP WITH TIME ZONE,
    duration INTERVAL, -- Calculated as left_at - joined_at
    recording_url TEXT,
    
    -- Unique constraint to prevent duplicate entries
    UNIQUE(room_id, user_id, joined_at)
);

-- Indexes
CREATE INDEX idx_meeting_history_room_id ON meeting_history(room_id);
CREATE INDEX idx_meeting_history_user_id ON meeting_history(user_id);
CREATE INDEX idx_meeting_history_joined_at ON meeting_history(joined_at);

-- Comments
COMMENT ON TABLE meeting_history IS 'Table to store meeting history';
COMMENT ON COLUMN meeting_history.id IS 'Unique identifier for the history entry';
COMMENT ON COLUMN meeting_history.room_id IS 'ID of the room (references rooms.id)';
COMMENT ON COLUMN meeting_history.user_id IS 'ID of the user (references users.id)';
COMMENT ON COLUMN meeting_history.joined_at IS 'Timestamp when user joined the meeting';
COMMENT ON COLUMN meeting_history.left_at IS 'Timestamp when user left the meeting';
COMMENT ON COLUMN meeting_history.duration IS 'Duration of the meeting participation';
COMMENT ON COLUMN meeting_history.recording_url IS 'URL to the meeting recording (if available)';
```

### 4.8 Tabel User Settings

Tabel `user_settings` menyimpan pengaturan pengguna.

```sql
CREATE TABLE user_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    setting_key VARCHAR(255) NOT NULL,
    setting_value TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Unique constraint to prevent duplicate settings
    UNIQUE(user_id, setting_key)
);

-- Indexes
CREATE INDEX idx_user_settings_user_id ON user_settings(user_id);
CREATE INDEX idx_user_settings_setting_key ON user_settings(setting_key);

-- Comments
COMMENT ON TABLE user_settings IS 'Table to store user settings';
COMMENT ON COLUMN user_settings.id IS 'Unique identifier for the setting';
COMMENT ON COLUMN user_settings.user_id IS 'ID of the user (references users.id)';
COMMENT ON COLUMN user_settings.setting_key IS 'Setting key';
COMMENT ON COLUMN user_settings.setting_value IS 'Setting value';
COMMENT ON COLUMN user_settings.created_at IS 'Timestamp when setting was created';
COMMENT ON COLUMN user_settings.updated_at IS 'Timestamp when setting was last updated';
```

### 4.9 Tabel Room Settings

Tabel `room_settings` menyimpan pengaturan ruang meeting.

```sql
CREATE TABLE room_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    setting_key VARCHAR(255) NOT NULL,
    setting_value TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Unique constraint to prevent duplicate settings
    UNIQUE(room_id, setting_key)
);

-- Indexes
CREATE INDEX idx_room_settings_room_id ON room_settings(room_id);
CREATE INDEX idx_room_settings_setting_key ON room_settings(setting_key);

-- Comments
COMMENT ON TABLE room_settings IS 'Table to store room settings';
COMMENT ON COLUMN room_settings.id IS 'Unique identifier for the setting';
COMMENT ON COLUMN room_settings.room_id IS 'ID of the room (references rooms.id)';
COMMENT ON COLUMN room_settings.setting_key IS 'Setting key';
COMMENT ON COLUMN room_settings.setting_value IS 'Setting value';
COMMENT ON COLUMN room_settings.created_at IS 'Timestamp when setting was created';
COMMENT ON COLUMN room_settings.updated_at IS 'Timestamp when setting was last updated';
```

### 4.10 Tabel Notifications

Tabel `notifications` menyimpan notifikasi untuk pengguna.

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('room_invitation', 'meeting_reminder', 'system', 'message')),
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- Comments
COMMENT ON TABLE notifications IS 'Table to store user notifications';
COMMENT ON COLUMN notifications.id IS 'Unique identifier for the notification';
COMMENT ON COLUMN notifications.user_id IS 'ID of the user (references users.id)';
COMMENT ON COLUMN notifications.type IS 'Notification type (room_invitation, meeting_reminder, system, message)';
COMMENT ON COLUMN notifications.title IS 'Notification title';
COMMENT ON COLUMN notifications.message IS 'Notification message';
COMMENT ON COLUMN notifications.is_read IS 'Flag indicating if notification is read';
COMMENT ON COLUMN notifications.created_at IS 'Timestamp when notification was created';
```

## 5. Database Functions dan Triggers

### 5.1 Update Timestamp Trigger

```sql
-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for tables with updated_at column
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_rooms_updated_at BEFORE UPDATE ON rooms
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_settings_updated_at BEFORE UPDATE ON user_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_room_settings_updated_at BEFORE UPDATE ON room_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 5.2 Calculate Meeting Duration Trigger

```sql
-- Function to calculate meeting duration
CREATE OR REPLACE FUNCTION calculate_meeting_duration()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.left_at IS NOT NULL AND NEW.joined_at IS NOT NULL THEN
        NEW.duration = NEW.left_at - NEW.joined_at;
    ELSE
        NEW.duration = NULL;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for meeting_history table
CREATE TRIGGER calculate_meeting_history_duration BEFORE INSERT OR UPDATE ON meeting_history
    FOR EACH ROW EXECUTE FUNCTION calculate_meeting_duration();
```

### 5.3 Soft Delete User Function

```sql
-- Function to soft delete user
CREATE OR REPLACE FUNCTION soft_delete_user(user_id UUID)
RETURNS BOOLEAN AS $$
DECLARE
    user_exists BOOLEAN;
BEGIN
    -- Check if user exists
    SELECT EXISTS(SELECT 1 FROM users WHERE id = user_id AND status != 'deleted') INTO user_exists;
    
    IF user_exists THEN
        -- Update user status to 'deleted'
        UPDATE users SET status = 'deleted', updated_at = NOW() WHERE id = user_id;
        
        -- Revoke all user sessions
        UPDATE user_sessions SET is_revoked = TRUE WHERE user_id = user_id;
        
        RETURN TRUE;
    ELSE
        RETURN FALSE;
    END IF;
END;
$$ language 'plpgsql';
```

### 5.4 Clean Expired Sessions Function

```sql
-- Function to clean expired sessions
CREATE OR REPLACE FUNCTION clean_expired_sessions()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    -- Delete expired and revoked sessions
    DELETE FROM user_sessions 
    WHERE expires_at < NOW() OR is_revoked = TRUE;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    
    RETURN deleted_count;
END;
$$ language 'plpgsql';
```

## 6. Database Views

### 6.1 Active Rooms View

```sql
-- View for active rooms
CREATE OR REPLACE VIEW active_rooms AS
SELECT 
    r.id,
    r.name,
    r.description,
    r.host_id,
    u.first_name || ' ' || u.last_name AS host_name,
    r.max_users,
    r.created_at,
    r.scheduled_at,
    COUNT(rp.user_id) AS current_participants
FROM 
    rooms r
    LEFT JOIN users u ON r.host_id = u.id
    LEFT JOIN room_participants rp ON r.id = rp.room_id AND rp.left_at IS NULL
WHERE 
    r.status = 'active'
GROUP BY 
    r.id, r.name, r.description, r.host_id, u.first_name, u.last_name, r.max_users, r.created_at, r.scheduled_at;

-- Comments
COMMENT ON VIEW active_rooms IS 'View showing active rooms with participant count';
```

### 6.2 User Contacts View

```sql
-- View for user contacts
CREATE OR REPLACE VIEW user_contacts_view AS
SELECT 
    uc.id,
    uc.user_id,
    uc.contact_id,
    u.first_name || ' ' || u.last_name AS contact_name,
    u.email AS contact_email,
    u.avatar AS contact_avatar,
    u.status AS contact_status,
    uc.created_at
FROM 
    user_contacts uc
    JOIN users u ON uc.contact_id = u.id
WHERE 
    uc.user_id IS NOT NULL AND u.status != 'deleted';

-- Comments
COMMENT ON VIEW user_contacts_view IS 'View showing user contacts with contact details';
```

### 6.3 Meeting Statistics View

```sql
-- View for meeting statistics
CREATE OR REPLACE VIEW meeting_statistics AS
SELECT 
    r.id AS room_id,
    r.name AS room_name,
    r.host_id,
    u.first_name || ' ' || u.last_name AS host_name,
    COUNT(DISTINCT mh.user_id) AS total_participants,
    COUNT(mh.id) AS total_participations,
    AVG(EXTRACT(EPOCH FROM mh.duration)) AS avg_duration_seconds,
    MIN(mh.joined_at) AS meeting_start,
    MAX(mh.left_at) AS meeting_end,
    r.created_at AS room_created_at,
    r.ended_at AS room_ended_at
FROM 
    rooms r
    LEFT JOIN users u ON r.host_id = u.id
    LEFT JOIN meeting_history mh ON r.id = mh.room_id
WHERE 
    r.status = 'ended'
GROUP BY 
    r.id, r.name, r.host_id, u.first_name, u.last_name, r.created_at, r.ended_at;

-- Comments
COMMENT ON VIEW meeting_statistics IS 'View showing meeting statistics';
```

### 6.4 User Activity View

```sql
-- View for user activity
CREATE OR REPLACE VIEW user_activity AS
SELECT 
    u.id AS user_id,
    u.first_name || ' ' || u.last_name AS user_name,
    u.email,
    COUNT(DISTINCT r.id) AS rooms_created,
    COUNT(DISTINCT rp.room_id) AS rooms_joined,
    COUNT(mh.id) AS meeting_participations,
    SUM(EXTRACT(EPOCH FROM mh.duration)) AS total_meeting_seconds,
    u.last_login,
    u.created_at
FROM 
    users u
    LEFT JOIN rooms r ON u.id = r.host_id
    LEFT JOIN room_participants rp ON u.id = rp.user_id
    LEFT JOIN meeting_history mh ON u.id = mh.user_id
WHERE 
    u.status != 'deleted'
GROUP BY 
    u.id, u.first_name, u.last_name, u.email, u.last_login, u.created_at;

-- Comments
COMMENT ON VIEW user_activity IS 'View showing user activity statistics';
```

## 7. Redis Data Structures

### 7.1 User Sessions

```redis
-- User session data
HSET user_sessions:{token} user_id {user_id}
HSET user_sessions:{token} expires_at {timestamp}
HSET user_sessions:{token} created_at {timestamp}
EXPIRE user_sessions:{token} 86400 -- 24 hours

-- User sessions index
SADD user_sessions:{user_id} {token}
```

### 7.2 Active Room Participants

```redis
-- Room participants
SADD room_participants:{room_id} {user_id}

-- User current room
SET user_current_room:{user_id} {room_id}
EXPIRE user_current_room:{user_id} 3600 -- 1 hour
```

### 7.3 WebRTC Signaling Data

```redis
-- WebRTC offers
HSET webrtc_offers:{room_id}:{user_id}:{target_user_id} {offer_data}
EXPIRE webrtc_offers:{room_id}:{user_id}:{target_user_id} 300 -- 5 minutes

-- WebRTC answers
HSET webrtc_answers:{room_id}:{user_id}:{target_user_id} {answer_data}
EXPIRE webrtc_answers:{room_id}:{user_id}:{target_user_id} 300 -- 5 minutes

-- WebRTC ICE candidates
LPUSH webrtc_ice_candidates:{room_id}:{user_id}:{target_user_id} {candidate_data}
EXPIRE webrtc_ice_candidates:{room_id}:{user_id}:{target_user_id} 300 -- 5 minutes
```

### 7.4 Rate Limiting

```redis
-- API rate limiting
INCR api_rate_limit:{user_id}:{endpoint}
EXPIRE api_rate_limit:{user_id}:{endpoint} 60 -- 1 minute

-- WebSocket rate limiting
INCR ws_rate_limit:{user_id}:{room_id}
EXPIRE ws_rate_limit:{user_id}:{room_id} 10 -- 10 seconds
```

### 7.5 Caching

```redis
-- User profile cache
HSET user_profile:{user_id} first_name {first_name}
HSET user_profile:{user_id} last_name {last_name}
HSET user_profile:{user_id} avatar {avatar_url}
HSET user_profile:{user_id} status {status}
EXPIRE user_profile:{user_id} 3600 -- 1 hour

-- Room info cache
HSET room_info:{room_id} name {room_name}
HSET room_info:{room_id} host_id {host_id}
HSET room_info:{room_id} max_users {max_users}
HSET room_info:{room_id} status {status}
EXPIRE room_info:{room_id} 1800 -- 30 minutes
```

## 8. Database Security

### 8.1 User Roles and Permissions

```sql
-- Create roles
CREATE ROLE app_user;
CREATE ROLE app_admin;

-- Grant permissions to roles
GRANT CONNECT ON DATABASE webrtc_meeting TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;

GRANT ALL PRIVILEGES ON DATABASE webrtc_meeting TO app_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO app_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO app_admin;

-- Set default permissions for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO app_user;
```

### 8.2 Row Level Security (RLS)

```sql
-- Enable Row Level Security
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE rooms ENABLE ROW LEVEL SECURITY;
ALTER TABLE room_participants ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_contacts ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE room_messages ENABLE ROW LEVEL SECURITY;
ALTER TABLE meeting_history ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_settings ENABLE ROW LEVEL SECURITY;
ALTER TABLE room_settings ENABLE ROW LEVEL SECURITY;
ALTER TABLE notifications ENABLE ROW LEVEL SECURITY;

-- Users can only view/update their own profile
CREATE POLICY user_isolation ON users FOR ALL USING (id = current_user_id());

-- Users can only view rooms they are participating in or they created
CREATE POLICY room_participant_isolation ON rooms FOR SELECT 
    USING (
        id IN (
            SELECT room_id FROM room_participants WHERE user_id = current_user_id()
        ) OR host_id = current_user_id()
    );

-- Users can only update rooms they created
CREATE POLICY room_host_update ON rooms FOR UPDATE 
    USING (host_id = current_user_id());

-- Users can only view room participants for rooms they are in
CREATE POLICY room_participant_visibility ON room_participants FOR SELECT 
    USING (
        room_id IN (
            SELECT room_id FROM room_participants WHERE user_id = current_user_id()
        )
    );

-- Users can only view their own contacts
CREATE POLICY user_contacts_isolation ON user_contacts FOR ALL 
    USING (user_id = current_user_id());

-- Users can only view/update their own sessions
CREATE POLICY user_session_isolation ON user_sessions FOR ALL 
    USING (user_id = current_user_id());

-- Users can only view messages for rooms they are in
CREATE POLICY room_message_visibility ON room_messages FOR SELECT 
    USING (
        room_id IN (
            SELECT room_id FROM room_participants WHERE user_id = current_user_id()
        )
    );

-- Users can only insert messages for rooms they are in
CREATE POLICY room_message_insert ON room_messages FOR INSERT 
    WITH CHECK (
        room_id IN (
            SELECT room_id FROM room_participants WHERE user_id = current_user_id()
        )
    );

-- Users can only view their own meeting history
CREATE POLICY meeting_history_isolation ON meeting_history FOR ALL 
    USING (user_id = current_user_id());

-- Users can only view/update their own settings
CREATE POLICY user_settings_isolation ON user_settings FOR ALL 
    USING (user_id = current_user_id());

-- Users can only view notifications for themselves
CREATE POLICY notification_isolation ON notifications FOR ALL 
    USING (user_id = current_user_id());
```

### 8.3 Data Encryption

```sql
-- Enable pgcrypto extension for encryption functions
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Function to encrypt sensitive data
CREATE OR REPLACE FUNCTION encrypt_data(data TEXT, secret TEXT)
RETURNS TEXT AS $$
BEGIN
    RETURN encode(encrypt(data::bytea, secret::bytea, 'aes'), 'hex');
END;
$$ LANGUAGE plpgsql;

-- Function to decrypt sensitive data
CREATE OR REPLACE FUNCTION decrypt_data(encrypted_data TEXT, secret TEXT)
RETURNS TEXT AS $$
BEGIN
    RETURN convert_from(decrypt(encrypted_data::bytea, secret::bytea, 'aes'), 'UTF8');
END;
$$ LANGUAGE plpgsql;
```

## 9. Database Migration Strategy

### 9.1 Migration Tooling

Menggunakan migration tool seperti **Goose** atau **Golang-Migrate** untuk mengelola skema database:

```bash
# Install Goose
go install github.com/pressly/goose/cmd/goose@latest

# Create migration
goose create add_user_table sql

# Run migration
goose -dir=migrations postgres "user=postgres password=postgres dbname=webrtc_meeting sslmode=disable" up

# Rollback migration
goose -dir=migrations postgres "user=postgres password=postgres dbname=webrtc_meeting sslmode=disable" down
```

### 9.2 Migration Files Structure

```
migrations/
├── 20230101000000_create_users_table.up.sql
├── 20230101000000_create_users_table.down.sql
├── 20230101000001_create_rooms_table.up.sql
├── 20230101000001_create_rooms_table.down.sql
├── 20230101000002_create_room_participants_table.up.sql
├── 20230101000002_create_room_participants_table.down.sql
├── 20230101000003_create_user_contacts_table.up.sql
├── 20230101000003_create_user_contacts_table.down.sql
├── 20230101000004_create_user_sessions_table.up.sql
├── 20230101000004_create_user_sessions_table.down.sql
├── 20230101000005_create_room_messages_table.up.sql
├── 20230101000005_create_room_messages_table.down.sql
├── 20230101000006_create_meeting_history_table.up.sql
├── 20230101000006_create_meeting_history_table.down.sql
├── 20230101000007_create_user_settings_table.up.sql
├── 20230101000007_create_user_settings_table.down.sql
├── 20230101000008_create_room_settings_table.up.sql
├── 20230101000008_create_room_settings_table.down.sql
├── 20230101000009_create_notifications_table.up.sql
├── 20230101000009_create_notifications_table.down.sql
├── 20230101000010_create_functions_and_triggers.up.sql
├── 20230101000010_create_functions_and_triggers.down.sql
├── 20230101000011_create_views.up.sql
├── 20230101000011_create_views.down.sql
├── 20230101000012_setup_security.up.sql
└── 20230101000012_setup_security.down.sql
```

### 9.3 Example Migration File

```sql
-- 20230101000000_create_users_table.up.sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    avatar TEXT,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'banned')),
    email_verified BOOLEAN DEFAULT FALSE,
    last_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);
```

```sql
-- 20230101000000_create_users_table.down.sql
DROP TABLE users;
```

## 10. Database Backup and Recovery

### 10.1 Backup Strategy

```bash
# Full backup
pg_dump -h localhost -U postgres -d webrtc_meeting -f /backups/webrtc_meeting_full_$(date +%Y%m%d_%H%M%S).sql

# Compressed backup
pg_dump -h localhost -U postgres -d webrtc_meeting | gzip > /backups/webrtc_meeting_full_$(date +%Y%m%d_%H%M%S).sql.gz

# Schema-only backup
pg_dump -h localhost -U postgres -d webrtc_meeting -s -f /backups/webrtc_meeting_schema_$(date +%Y%m%d_%H%M%S).sql

# Data-only backup
pg_dump -h localhost -U postgres -d webrtc_meeting -a -f /backups/webrtc_meeting_data_$(date +%Y%m%d_%H%M%S).sql
```

### 10.2 Automated Backup Script

```bash
#!/bin/bash

# Database backup script
BACKUP_DIR="/backups"
DB_NAME="webrtc_meeting"
DB_USER="postgres"
DB_HOST="localhost"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Create full backup
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME | gzip > $BACKUP_DIR/${DB_NAME}_full_${TIMESTAMP}.sql.gz

# Remove backups older than RETENTION_DAYS
find $BACKUP_DIR -name "${DB_NAME}_full_*.sql.gz" -mtime +$RETENTION_DAYS -delete

# Log backup
echo "Backup created: ${DB_NAME}_full_${TIMESTAMP}.sql.gz" >> $BACKUP_DIR/backup_log.txt
```

### 10.3 Recovery Process

```bash
# Restore from full backup
gunzip -c /backups/webrtc_meeting_full_20230101_120000.sql.gz | psql -h localhost -U postgres -d webrtc_meeting

# Restore from schema-only backup
psql -h localhost -U postgres -d webrtc_meeting -f /backups/webrtc_meeting_schema_20230101_120000.sql

# Restore from data-only backup
psql -h localhost -U postgres -d webrtc_meeting -f /backups/webrtc_meeting_data_20230101_120000.sql
```

## 11. Database Monitoring and Performance

### 11.1 Monitoring Queries

```sql
-- Top 10 most expensive queries
SELECT query, calls, total_time, mean_time, rows
FROM pg_stat_statements
ORDER BY total_time DESC
LIMIT 10;

-- Table sizes
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Index usage
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;

-- Active connections
SELECT 
    datname,
    usename,
    application_name,
    client_addr,
    state,
    state_change,
    query
FROM pg_stat_activity
WHERE state != 'idle';
```

### 11.2 Performance Optimization

```sql
-- Update statistics for better query planning
ANALYZE;

-- Reindex frequently updated tables
REINDEX TABLE room_participants;
REINDEX TABLE room_messages;
REINDEX TABLE user_sessions;

-- Vacuum to reclaim storage
VACUUM;
VACUUM ANALYZE;

-- Set work memory for complex queries
SET work_mem = '64MB';

-- Set shared buffers for better caching
-- This should be set in postgresql.conf
-- shared_buffers = 256MB
```

### 11.3 Connection Pooling

Menggunakan **PgBouncer** untuk connection pooling:

```ini
; pgbouncer.ini
[databases]
webrtc_meeting = host=localhost port=5432 dbname=webrtc_meeting

[pgbouncer]
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 20
reserve_pool = 5
reserve_pool_timeout = 3
server_idle_timeout = 30
```

## 12. Kesimpulan

Desain database untuk aplikasi WebRTC meeting ini dirancang dengan mempertimbangkan:

1. **Struktur yang terorganisir** dengan tabel-tabel yang terpisah berdasarkan fungsinya
2. **Relasi yang jelas** antara entitas-entitas utama seperti users, rooms, dan participants
3. **Performa yang optimal** dengan indexing yang tepat dan query optimization
4. **Keamanan data** dengan encryption, row level security, dan proper user roles
5. **Scalability** dengan desain yang dapat menangani pertumbuhan data dan pengguna
6. **Maintainability** dengan migration strategy, backup and recovery plan
7. **Monitoring** dengan queries untuk memantau performa dan kesehatan database

Dengan desain database ini, aplikasi WebRTC meeting akan memiliki foundation yang kuat untuk mendukung fitur-fitur utama seperti user management, room-based meeting system, dan riwayat meeting.