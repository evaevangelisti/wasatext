<template>
  <div class="profile">
    <label class="profile__picture-field">
      <img
        :src="
          user?.profilePicture
            ? backendBaseUrl + user.profilePicture
            : defaultProfilePicture
        "
        alt="Profile picture"
        class="profile__picture"
        :class="{ 'profile__picture--error': photoError }"
      >
      <svg viewBox="0 0 24 24" fill="none" class="profile__picture-icon">
        <path
          d="M20.1498 7.93997L8.27978 19.81C7.21978 20.88 4.04977 21.3699 3.32977 20.6599C2.60977 19.9499 3.11978 16.78 4.17978 15.71L16.0498 3.84C16.5979 3.31801 17.3283 3.03097 18.0851 3.04019C18.842 3.04942 19.5652 3.35418 20.1004 3.88938C20.6356 4.42457 20.9403 5.14781 20.9496 5.90463C20.9588 6.66146 20.6718 7.39189 20.1498 7.93997V7.93997Z"
          stroke="var(--color-secondary)"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
      <input
        type="file"
        accept="image/*"
        style="display: none"
        @change="onPhotoChange"
      >
    </label>
    <div class="profile__username-field">
      <span class="text-secondary">Username</span>
      <div>
        <div
          class="profile__username-input"
          :class="{ 'profile__username-input--error': usernameError }"
        >
          <input
            ref="usernameInput"
            v-model="editedUsername"
            class="profile__username"
            :disabled="!editing"
          >
          <button
            ref="usernameButton"
            class="profile__username-icon-button"
            @click="editing ? saveUsername() : startEditing()"
          >
            <svg
              v-if="!editing"
              viewBox="0 0 24 24"
              fill="none"
              class="profile__username-icon"
            >
              <path
                d="M20.1498 7.93997L8.27978 19.81C7.21978 20.88 4.04977 21.3699 3.32977 20.6599C2.60977 19.9499 3.11978 16.78 4.17978 15.71L16.0498 3.84C16.5979 3.31801 17.3283 3.03097 18.0851 3.04019C18.842 3.04942 19.5652 3.35418 20.1004 3.88938C20.6356 4.42457 20.9403 5.14781 20.9496 5.90463C20.9588 6.66146 20.6718 7.39189 20.1498 7.93997V7.93997Z"
                stroke="var(--color-tertiary)"
                stroke-width="1.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
            <svg
              v-else
              viewBox="0 0 24 24"
              fill="none"
              class="profile__username-icon"
            >
              <path
                d="M4 12.6111L8.92308 17.5L20 6.5"
                stroke="var(--color-secondary)"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </button>
        </div>
        <span
          v-if="usernameError"
          class="username-error-message"
        >
          Invalid or unavailable username
        </span>
      </div>
    </div>
    <button class="logout-button" @click="logout">
      <span class="text-body">Logout</span>
      <svg width="36px" height="36px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path
          d="M12.9999 2C10.2385 2 7.99991 4.23858 7.99991 7C7.99991 7.55228 8.44762 8 8.99991 8C9.55219 8 9.99991 7.55228 9.99991 7C9.99991 5.34315 11.3431 4 12.9999 4H16.9999C18.6568 4 19.9999 5.34315 19.9999 7V17C19.9999 18.6569 18.6568 20 16.9999 20H12.9999C11.3431 20 9.99991 18.6569 9.99991 17C9.99991 16.4477 9.55219 16 8.99991 16C8.44762 16 7.99991 16.4477 7.99991 17C7.99991 19.7614 10.2385 22 12.9999 22H16.9999C19.7613 22 21.9999 19.7614 21.9999 17V7C21.9999 4.23858 19.7613 2 16.9999 2H12.9999Z"
          fill="var(--color-error)"
        />
        <path
          d="M13.9999 11C14.5522 11 14.9999 11.4477 14.9999 12C14.9999 12.5523 14.5522 13 13.9999 13V11Z"
          fill="var(--color-error)"
        />
        <path
          d="M5.71783 11C5.80685 10.8902 5.89214 10.7837 5.97282 10.682C6.21831 10.3723 6.42615 10.1004 6.57291 9.90549C6.64636 9.80795 6.70468 9.72946 6.74495 9.67492L6.79152 9.61162L6.804 9.59454L6.80842 9.58848C6.80846 9.58842 6.80892 9.58778 5.99991 9L6.80842 9.58848C7.13304 9.14167 7.0345 8.51561 6.58769 8.19098C6.14091 7.86637 5.51558 7.9654 5.19094 8.41215L5.18812 8.41602L5.17788 8.43002L5.13612 8.48679C5.09918 8.53682 5.04456 8.61033 4.97516 8.7025C4.83623 8.88702 4.63874 9.14542 4.40567 9.43937C3.93443 10.0337 3.33759 10.7481 2.7928 11.2929L2.08569 12L2.7928 12.7071C3.33759 13.2519 3.93443 13.9663 4.40567 14.5606C4.63874 14.8546 4.83623 15.113 4.97516 15.2975C5.04456 15.3897 5.09918 15.4632 5.13612 15.5132L5.17788 15.57L5.18812 15.584L5.19045 15.5872C5.51509 16.0339 6.14091 16.1336 6.58769 15.809C7.0345 15.4844 7.13355 14.859 6.80892 14.4122L5.99991 15C6.80892 14.4122 6.80897 14.4123 6.80892 14.4122L6.804 14.4055L6.79152 14.3884L6.74495 14.3251C6.70468 14.2705 6.64636 14.1921 6.57291 14.0945C6.42615 13.8996 6.21831 13.6277 5.97282 13.318C5.89214 13.2163 5.80685 13.1098 5.71783 13H13.9999V11H5.71783Z"
          fill="var(--color-error)"
        />
      </svg>
    </button>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from "vue";
import api from "@/services/api";
import { backendBaseUrl } from "@/services/baseUrl";

import defaultProfilePicture from "@/assets/default-profile-picture.jpg";

const props = defineProps({
  user: Object,
});

const emit = defineEmits(["profile-update", "logout"]);

const editing = ref(false);
const editedUsername = ref(props.user?.username || "");
const usernameInput = ref(null);
const usernameButton = ref(null);
const usernameError = ref(false);

watch(
  () => props.user?.username,
  (newValue) => {
    editedUsername.value = newValue || "";
    editing.value = false;
  },
);

function startEditing() {
  editing.value = true;
  nextTick(() => {
    usernameInput.value?.focus();
  });
}

async function saveUsername() {
  if (editedUsername.value === props.user?.username) {
    usernameInput.value?.blur();
    usernameButton.value?.blur();
    editing.value = false;
    return;
  }

  usernameError.value = false;

  try {
    const response = await api.put("/me/username", {
      username: editedUsername.value,
    });

    emit("profile-update", response.data);
    editing.value = false;

    usernameInput.value?.blur();
    usernameButton.value?.blur();
  } catch (e) {
    console.error(e);
    usernameError.value = true;

    usernameInput.value?.blur();
    usernameButton.value?.blur();
  }
}

const photoError = ref(false);

async function onPhotoChange(event) {
  const file = event.target.files[0];
  if (!file) return;

  const formData = new FormData();
  formData.append("image", file);

  photoError.value = false;

  try {
    const response = await api.put("/me/photo", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });

    emit("profile-update", response.data);
  } catch (e) {
    photoError.value = true;
  }
}

function logout() {
  delete api.defaults.headers.common["Authorization"];
  emit("logout");
}
</script>

<style>
.profile {
  display: flex;
  align-items: center;
  flex-direction: column;
  gap: 2.5rem;
  padding: 3rem;
}

.profile__picture-field {
  position: relative;
  display: inline-block;
  cursor: pointer;
}

.profile__picture {
  border: 2px solid transparent;
  border-radius: 50%;
  width: 128px;
  height: 128px;
  object-fit: cover;
  transition: opacity 0.1s;
}

.profile__picture--error {
  border-color: var(--color-error) !important;
}

.profile__picture-field:hover .profile__picture {
  opacity: 0.5;
}

.profile__picture-icon {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  pointer-events: none;
  width: 24px;
  height: 24px;
  opacity: 0;
  transition: opacity 0.1s;
}

.profile__picture-field:hover .profile__picture-icon {
  opacity: 1;
}

.profile__username-field {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  width: 100%;
}

.profile__username-input {
  display: flex;
  border-bottom: 2px solid var(--color-tertiary);
}

.profile__username-input--error {
  border-bottom-color: var(--color-error);
}

.profile__username-input:focus-within {
  border-bottom: 2px solid var(--color-primary);
}

.profile__username {
  border: none;
  width: 100%;
  padding: 0.5rem 0rem;
  background-color: inherit;
  font-size: 1rem;
  color: var(--color-secondary);
}

.profile__username-icon-button {
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  padding: 0.2rem;
  background-color: inherit;
}

.profile__username-icon {
  width: 24px;
  height: 24px;
}

.username-error-message {
  color: var(--color-error);
  font-size: 0.9em;
  margin-top: 0.25em;
}

.logout-button {
  align-self: flex-end;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border: none;
  padding: 0;
  background-color: inherit;
  color: var(--color-error);
}

.logout-button .text-body,
.logout-button svg {
  transition: filter 0.1s;
}

.logout-button:hover svg,
.logout-button:focus svg {
  filter: brightness(1.1);
}

.logout-button:hover .text-body,
.logout-button:focus .text-body {
  filter: brightness(1.1);
}
</style>
