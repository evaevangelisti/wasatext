<template>
  <div class="conversations-wrapper">
    <ul class="conversations">
      <li
        v-for="conversation in conversations"
        :key="conversation.conversationId"
        class="conversation-item"
      >
        <button
          class="conversation__button"
          :class="{ active: isActive(conversation) }"
          @click="selectConversation(conversation)"
        >
          <img
            v-if="conversation.type === 'group'"
            :src="
              conversation.photo
                ? backendBaseUrl + conversation.photo
                : defaultGroupPicture
            "
            alt="Group Picture"
            class="conversation-photo"
          >
          <img
            v-else
            :src="
              getOtherUser(conversation)?.profilePicture
                ? backendBaseUrl + getOtherUser(conversation).profilePicture
                : defaultProfilePicture
            "
            alt="Profile Picture"
            class="conversation-photo"
          >
          <div class="conversation__info-wrapper">
            <div class="conversation__info">
              <span class="text-body" style="font-weight: 600">
                {{
                  conversation.type === "group"
                    ? conversation.name
                    : getOtherUser(conversation)?.username
                }}
              </span>
              <span class="text-caption">{{
                conversation.lastMessage &&
                  formatSentAt(conversation.lastMessage.sentAt)
              }}</span>
            </div>
            <div class="conversation__preview">
              <template v-if="conversation.lastMessage">
                <template v-if="conversation.type === 'group'">
                  <span
                    class="conversation__sender text-body"
                    style="color: var(--color-tertiary)"
                  >
                    {{
                      conversation.lastMessage.sender?.userId === user.userId
                        ? "You"
                        : conversation.lastMessage.sender?.username
                    }}<span>:&nbsp;</span>
                  </span>
                </template>
                <template
                  v-if="
                    conversation.lastMessage &&
                      conversation.lastMessage.attachment
                  "
                >
                  <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    class="conversation__attachment-icon"
                  >
                    <path
                      d="M14.2639 15.9375L12.5958 14.2834C11.7909 13.4851 11.3884 13.086 10.9266 12.9401C10.5204 12.8118 10.0838 12.8165 9.68048 12.9536C9.22188 13.1095 8.82814 13.5172 8.04068 14.3326L4.04409 18.2801M14.2639 15.9375L14.6053 15.599C15.4112 14.7998 15.8141 14.4002 16.2765 14.2543C16.6831 14.126 17.12 14.1311 17.5236 14.2687C17.9824 14.4251 18.3761 14.8339 19.1634 15.6514L20 16.4934M14.2639 15.9375L18.275 19.9565M18.275 19.9565C17.9176 20 17.4543 20 16.8 20H7.2C6.07989 20 5.51984 20 5.09202 19.782C4.71569 19.5903 4.40973 19.2843 4.21799 18.908C4.12796 18.7313 4.07512 18.5321 4.04409 18.2801M18.275 19.9565C18.5293 19.9256 18.7301 19.8727 18.908 19.782C19.2843 19.5903 19.5903 19.2843 19.782 18.908C20 18.4802 20 17.9201 20 16.8V16.4934M4.04409 18.2801C4 17.9221 4 17.4575 4 16.8V7.2C4 6.0799 4 5.51984 4.21799 5.09202C4.40973 4.71569 4.71569 4.40973 5.09202 4.21799C5.51984 4 6.07989 4 7.2 4H16.8C17.9201 4 18.4802 4 18.908 4.21799C19.2843 4.40973 19.5903 4.71569 19.782 5.09202C20 5.51984 20 6.0799 20 7.2V16.4934M17 8.99989C17 10.1045 16.1046 10.9999 15 10.9999C13.8954 10.9999 13 10.1045 13 8.99989C13 7.89532 13.8954 6.99989 15 6.99989C16.1046 6.99989 17 7.89532 17 8.99989Z"
                      stroke="var(--color-tertiary)"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                  <span>&nbsp;</span>
                  <span
                    class="conversation__content text-body"
                    style="color: var(--color-tertiary)"
                  >{{ conversation.lastMessage.content }}</span>
                </template>
                <template v-else>
                  <span
                    class="conversation__content text-body"
                    style="color: var(--color-tertiary)"
                  >{{ conversation.lastMessage.content }}</span>
                </template>
              </template>
              <template v-else>
                <span class="conversation__content text-secondary">&nbsp;</span>
              </template>
            </div>
          </div>
        </button>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from "vue";
import api from "@/services/api";
import { backendBaseUrl } from "@/services/baseUrl";

import defaultProfilePicture from "@/assets/default-profile-picture.jpg";
import defaultGroupPicture from "@/assets/default-group-picture.jpg";

const props = defineProps({
  user: Object,
  activeConversation: Object,
  conversationUpdated: [Number, String, Boolean],
});

const emit = defineEmits(["active-conversation"]);

const conversations = ref([]);

function getOtherUser(conversation) {
  if (conversation.type !== "private") return null;
  return conversation.participants.find(
    (user) => user.userId !== props.user?.userId,
  );
}

async function loadConversations() {
  try {
    const response = await api.get("/conversations");
    let filtered = response.data.filter(
      (c) => (c.type === "private" && c.lastMessage) || c.type === "group",
    );

    conversations.value = filtered;
  } catch (e) {
    conversations.value = [];
  }
}

onMounted(loadConversations);

watch(
  () => props.conversationUpdated,
  () => {
    loadConversations();
  },
);

function selectConversation(conversation) {
  emit("active-conversation", conversation);
}

function isActive(conversation) {
  return (
    props.activeConversation &&
    conversation.conversationId === props.activeConversation.conversationId
  );
}

function formatSentAt(sentAt) {
  if (!sentAt) return "";
  const date = new Date(sentAt);
  const now = new Date();

  if (
    date.getFullYear() === now.getFullYear() &&
    date.getMonth() === now.getMonth() &&
    date.getDate() === now.getDate()
  ) {
    return date.toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
    });
  }

  const yesterday = new Date(now);
  yesterday.setDate(now.getDate() - 1);
  if (
    date.getFullYear() === yesterday.getFullYear() &&
    date.getMonth() === yesterday.getMonth() &&
    date.getDate() === yesterday.getDate()
  ) {
    return "Yesterday";
  }

  return date.toLocaleDateString([], {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
}
</script>

<style>
.conversations-wrapper {
  flex: 1 1 auto;
  overflow-y: auto;
  min-height: 0;
}

.conversations {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  width: 100%;
  padding: 1rem 1.5rem;
}

.conversation-item {
  display: flex;
  width: 100%;
}

.conversation__button {
  display: flex;
  align-items: center;
  gap: 1rem;
  border: none;
  border-radius: 12px;
  width: 100%;
  padding: 0.5rem 0.75rem;
  background-color: inherit;
  color: var(--color-secondary);
  transition: background-color 0.1s;
  overflow: hidden;
}

.conversation__button:hover,
.conversation__button:focus,
.conversation__button.active {
  outline: none;
  background-color: var(--color-quaternary);
}

.conversation-photo {
  width: 48px !important;
  height: 48px !important;
  aspect-ratio: 1 / 1;
  border: 2px transparent;
  border-radius: 50%;
  object-fit: cover;
}

.conversation__info-wrapper {
  display: flex;
  flex-direction: column;
  width: 100%;
  min-width: 0;
  overflow: hidden;
}

.conversation__info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.conversation__preview {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: flex-start;
}

.conversation__attachment-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  flex-grow: 0;
  flex-basis: 20px;
  min-width: 20px;
  min-height: 20px;
  max-width: 20px;
  max-height: 20px;
  display: inline-block;
  vertical-align: middle;
}

.conversation__content {
  display: block;
  max-width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
