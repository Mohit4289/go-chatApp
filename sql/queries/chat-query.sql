-- name: CreateConversation :one
INSERT INTO public.conversations (type, title)
VALUES ($1, $2)
RETURNING id, type, title, created_at, updated_at;

-- name: AddParticipant :one
INSERT INTO public.conversation_participants (conversation_id, user_id, role)
VALUES ($1, $2, $3)
RETURNING id, conversation_id, user_id, role, joined_at, last_read_at;

-- name: GetUserConversations :many
SELECT 
    c.id AS conversation_id,
    c.type,
    c.title,
    c.updated_at,
    (
        SELECT m.content FROM public.messages m 
        WHERE m.conversation_id = c.id 
        ORDER BY m.created_at DESC LIMIT 1
    ) AS last_message,
    (
        SELECT m.created_at FROM public.messages m 
        WHERE m.conversation_id = c.id 
        ORDER BY m.created_at DESC LIMIT 1
    ) AS last_message_time
FROM public.conversations c
JOIN public.conversation_participants cp ON cp.conversation_id = c.id
WHERE cp.user_id = $1
ORDER BY c.updated_at DESC;

-- name: CreateMessage :one
INSERT INTO public.messages (conversation_id, sender_id, content, message_type)
VALUES ($1, $2, $3, $4)
RETURNING id, conversation_id, sender_id, content, message_type, created_at;

-- name: GetMessagesByConversation :many
SELECT 
    m.id,
    m.conversation_id,
    m.sender_id,
    u.name AS sender_name,
    m.content,
    m.message_type,
    m.created_at
FROM public.messages m
JOIN public."user" u ON u.id = m.sender_id
WHERE m.conversation_id = $1
ORDER BY m.created_at ASC
LIMIT $2 OFFSET $3;

-- name: IsUserInConversation :one
SELECT EXISTS (
    SELECT 1 FROM public.conversation_participants
    WHERE conversation_id = $1 AND user_id = $2
);

-- name: GetConversationParticipants :many

SELECT user_id
FROM public.conversation_participants
WHERE conversation_id = $1;