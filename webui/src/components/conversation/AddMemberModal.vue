<template>
  <div class="add-member-modal-overlay" @click.self="close">
    <div class="add-member-modal">
      <div class="add-member-modal__header">
        <span class="text-body">Add member</span>
        <button
          aria-label="Close"
          class="add-member-modal__close"
          @click="close"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            class="add-member-modal__close__icon"
          >
            <path
              d="M20.7457 3.32851C20.3552 2.93798 19.722 2.93798 19.3315 3.32851L12.0371 10.6229L4.74275 3.32851C4.35223 2.93798 3.71906 2.93798 3.32854 3.32851C2.93801 3.71903 2.93801 4.3522 3.32854 4.74272L10.6229 12.0371L3.32856 19.3314C2.93803 19.722 2.93803 20.3551 3.32856 20.7457C3.71908 21.1362 4.35225 21.1362 4.74277 20.7457L12.0371 13.4513L19.3315 20.7457C19.722 21.1362 20.3552 21.1362 20.7457 20.7457C21.1362 20.3551 21.1362 19.722 20.7457 19.3315L13.4513 12.0371L20.7457 4.74272C21.1362 4.3522 21.1362 3.71903 20.7457 3.32851Z"
              fill="var(--color-tertiary)"
            />
          </svg>
        </button>
      </div>
      <div class="query">
        <svg viewBox="0 0 24 24" fill="none" class="query-icon">
          <path
            fill-rule="evenodd"
            clip-rule="evenodd"
            d="M15 10.5C15 12.9853 12.9853 15 10.5 15C8.01472 15 6 12.9853 6 10.5C6 8.01472 8.01472 6 10.5 6C12.9853 6 15 8.01472 15 10.5ZM14.1793 15.2399C13.1632 16.0297 11.8865 16.5 10.5 16.5C7.18629 16.5 4.5 13.8137 4.5 10.5C4.5 7.18629 7.18629 4.5 10.5 4.5C13.8137 4.5 16.5 7.18629 16.5 10.5C16.5 11.8865 16.0297 13.1632 15.2399 14.1792L20.0304 18.9697L18.9697 20.0303L14.1793 15.2399Z"
            fill="var(--color-tertiary)"
          />
        </svg>
        <input
          v-model="query"
          type="text"
          placeholder="Search"
          class="query-input"
        >
      </div>
      <div class="users-wrapper">
        <ul class="users">
          <li v-for="user in usersToAdd" :key="user.userId" class="user-item">
            <button class="user__button" @click="addMember(user)">
              <img
                :src="resolveImageUrl(user.profilePicture, defaultProfilePicture)"
                alt="Profile Picture"
                class="user__picture"
              >
              <span class="text-body" style="font-weight: 600">{{
                user.username
              }}</span>
            </button>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from "vue";
import api from "@/services/api";
import { resolveImageUrl } from "@/services/imageUrl";
import defaultProfilePicture from "@/assets/default-profile-picture.jpg";

const props = defineProps({
  conversation: Object,
  show: Boolean,
});

const emit = defineEmits(["close", "member-added"]);

const allUsers = ref([]);
const query = ref("");

async function loadUsers() {
  const response = await api.get("/users");
  allUsers.value = response.data;
}

loadUsers();

watch(query, async (newQuery) => {
  if (!newQuery.trim()) {
    loadUsers();
    return;
  }
  try {
    const response = await api.get("/users", {
      params: { q: newQuery },
    });
    allUsers.value = response.data;
  } catch (e) {
    allUsers.value = [];
  }
});

const usersToAdd = computed(() =>
  allUsers.value.filter(
    (u) => !props.conversation.members.some((m) => m.userId === u.userId),
  ),
);

async function addMember(user) {
  try {
    const response = await api.post(
      `/groups/${props.conversation.conversationId}/members`,
      { userId: user.userId },
    );

    emit("member-added", response.data);
    close();
  } catch (e) {
    console.error(e);
  }
}

function close() {
  emit("close");
}
</script>

<style scoped>
.add-member-modal-overlay {
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

.add-member-modal {
  background: var(--color-background);
  border-radius: 16px;
  padding: 1rem 0rem;
  min-width: 400px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.add-member-modal__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0rem 1rem;
}

.add-member-modal__close {
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 50%;
  padding: 0.3rem;
  background-color: inherit;
}

.add-member-modal__close:hover {
  background-color: var(--color-quaternary);
  transition: background 0.1s;
}

.add-member-modal__close__icon {
  width: 20px;
  height: 20px;
}

.add-member-modal > .query {
  margin: 0rem 1rem;
}

.users-wrapper {
  max-height: 60vh;
  overflow-y: auto;
}

.users {
  padding: 0rem 1rem;
}
</style>
