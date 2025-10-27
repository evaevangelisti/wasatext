<template>
  <div class="forward-modal-overlay" @click.self="close">
    <div class="forward-modal">
      <div class="forward-modal__header">
        <span class="text-body">Forward message to</span>
        <button aria-label="Close" class="forward-modal__close" @click="close">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            class="forward-modal__close__icon"
          >
            <path
              d="M20.7457 3.32851C20.3552 2.93798 19.722 2.93798 19.3315 3.32851L12.0371 10.6229L4.74275 3.32851C4.35223 2.93798 3.71906 2.93798 3.32854 3.32851C2.93801 3.71903 2.93801 4.3522 3.32854 4.74272L10.6229 12.0371L3.32856 19.3314C2.93803 19.722 2.93803 20.3551 3.32856 20.7457C3.71908 21.1362 4.35225 21.1362 4.74277 20.7457L12.0371 13.4513L19.3315 20.7457C19.722 21.1362 20.3552 21.1362 20.7457 20.7457C21.1362 20.3551 21.1362 19.722 20.7457 19.3315L13.4513 12.0371L20.7457 4.74272C21.1362 4.3522 21.1362 3.71903 20.7457 3.32851Z"
              fill="var(--color-tertiary)"
            />
          </svg>
        </button>
      </div>
      <div class="conversations-wrapper">
        <ul class="conversations">
          <li
            v-for="item in forwardList"
            :key="item.type === 'conversation' ? item.data.conversationId : item.data.userId"
            class="conversation-item"
          >
            <button
              class="conversation__button"
              @click="item.type === 'conversation'
                ? selectConversation(item.data)
                : forwardToUser(item.data)"
            >
              <template v-if="item.type === 'conversation'">
                <img
                  v-if="item.data.type === 'group'"
                  :src="resolveImageUrl(item.data.photo, defaultGroupPicture)"
                  alt="Group Picture"
                  class="conversation-photo"
                >
                <img
                  v-else
                  :src="resolveImageUrl(getOtherUser(item.data)?.profilePicture, defaultProfilePicture)"
                  alt="Profile Picture"
                  class="conversation-photo"
                >
                <span class="text-body" style="font-weight: 600">
                  {{
                    item.data.type === "group"
                      ? item.data.name
                      : getOtherUser(item.data)?.username
                  }}
                </span>
              </template>
              <template v-else>
                <img
                  :src="resolveImageUrl(item.data.profilePicture, defaultProfilePicture)"
                  alt="Profile Picture"
                  class="conversation-photo"
                >
                <span class="text-body" style="font-weight: 600">
                  {{ item.data.username }}
                </span>
              </template>
            </button>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import api from "@/services/api";
import { resolveImageUrl } from "@/services/imageUrl";
import defaultProfilePicture from "@/assets/default-profile-picture.jpg";
import defaultGroupPicture from "@/assets/default-group-picture.jpg";

const props = defineProps({
  user: Object,
  show: Boolean,
});

const emit = defineEmits(["close", "forward"]);

const conversations = ref([]);
const allUsers = ref([]);

function getOtherUser(conversation) {
  if (conversation.type !== "private") return null;
  return conversation.participants.find(
    (user) => user.userId !== props.user?.userId,
  );
}

async function loadConversations() {
  const response = await api.get("/conversations");
  conversations.value = response.data.filter(
    (c) => (c.type === "private" && c.lastMessage) || c.type === "group",
  );
}

async function loadUsers() {
  const response = await api.get("/users");
  allUsers.value = response.data.filter(u => u.userId !== props.user.userId);
}

onMounted(() => {
  loadConversations();
  loadUsers();
});

const forwardList = computed(() => {
  const privateUserIds = conversations.value
    .filter(c => c.type === "private")
    .map(c => getOtherUser(c)?.userId);

  const convItems = conversations.value.map(c => ({
    type: "conversation",
    data: c,
  }));

  const userItems = allUsers.value
    .filter(u => !privateUserIds.includes(u.userId))
    .map(u => ({
      type: "user",
      data: u,
    }));

  return [...convItems, ...userItems];
});

function selectConversation(conversation) {
  emit("forward", conversation);
}

async function forwardToUser(user) {
  try {
    const response = await api.post("/conversations", {
      type: "private",
      userId: user.userId,
    });
    emit("forward", response.data);
  } catch (e) {
    if (e.response && e.response.status === 409) {
      emit("forward", e.response.data);
    }
  }
}

function close() {
  emit("close");
}
</script>

<style scoped>
.forward-modal-overlay {
  position: fixed;
  z-index: 10001;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
}

.forward-modal {
  background: var(--color-background);
  border-radius: 16px;
  padding: 1rem 0rem;
  min-width: 400px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.forward-modal__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0rem 1rem;
}

.forward-modal__close {
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 50%;
  padding: 0.3rem;
  background-color: inherit;
}

.forward-modal__close:hover {
  background-color: var(--color-quaternary);
  transition: background 0.1s;
}

.forward-modal__close__icon {
  width: 20px;
  height: 20px;
}

.conversations-wrapper {
  max-height: 60vh;
  overflow-y: auto;
}

.conversations {
  padding: 0rem 1rem;
}
</style>
