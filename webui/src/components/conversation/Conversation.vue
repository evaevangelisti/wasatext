<template>
  <div
    v-if="props.conversation"
    class="conversation"
    :class="{
      'conversation--sidebar-open': sidebarOpen,
      'conversation--modal-open': forwardModalOpen,
    }"
  >
    <div class="conversation-header">
      <button
        class="conversation-header__button"
        @click="$emit('show-info-sidebar')"
      >
        <img
          v-if="conversation.type === 'group'"
          :src="
            conversation.photo
              ? backendBaseUrl + conversation.photo
              : defaultGroupPicture
          "
          alt="Group Picture"
          class="conversation-header__photo"
        >
        <img
          v-else
          :src="
            getOtherUser(conversation)?.profilePicture
              ? backendBaseUrl + getOtherUser(conversation).profilePicture
              : defaultProfilePicture
          "
          alt="Profile Picture"
          class="conversation-header__photo"
        >
        <span class="text-body" style="font-weight: 600">
          {{
            conversation.type === "group"
              ? conversation.name
              : getOtherUser(conversation)?.username
          }}
        </span>
      </button>
    </div>
    <div ref="messagesContainer" class="messages-wrapper">
      <div class="messages">
        <template v-for="(msg, idx) in messages" :key="msg.messageId">
          <div
            v-if="
              idx === 0 || !isSameDay(msg.sentAt, messages[idx - 1].sentAt)
            "
            class="date-separator"
          >
            <span class="text-secondary">{{ formatDay(msg.sentAt) }}</span>
          </div>
          <div
            class="message-wrapper"
            :class="{
              'message-wrapper--mine': msg.sender.userId === user.userId,
            }"
          >
            <template
              v-if="
                conversation.type === 'group' &&
                  msg.sender.userId !== user.userId &&
                  (
                    idx === 0 ||
                    messages[idx - 1].sender.userId !== msg.sender.userId ||
                    !isSameDay(msg.sentAt, messages[idx - 1].sentAt)
                  )
              "
            >
              <img
                class="message__avatar"
                :src="
                  msg.sender.profilePicture
                    ? backendBaseUrl + msg.sender.profilePicture
                    : defaultProfilePicture
                "
                alt="Profile"
              >
            </template>
            <div
              :ref="(el) => setMessageRef(msg.messageId, el)"
              class="message"
              :class="{
                'message--mine': msg.sender.userId === user.userId,
                editing: editingMessageId === msg.messageId,
              }"
              @mouseenter="hoveredMessageId = msg.messageId"
              @mouseleave="hoveredMessageId = null"
            >
              <template
                v-if="
                  conversation.type === 'group' &&
                    msg.sender.userId !== user.userId
                "
              >
                <span class="text-body sender-name" style="font-weight: 600">{{
                  msg.sender.username
                }}</span>
              </template>
              <div
                v-if="editingMessageId === msg.messageId"
                class="message__content"
              >
                <textarea
                  ref="editInput"
                  v-model="editingContent"
                  rows="1"
                  class="message__edit"
                  @keydown.enter="onEditEnter($event, msg)"
                  @keydown.esc="cancelEdit"
                  @input="autoResizeEdit"
                />
              </div>
              <div
                v-else
                class="message__content"
                :class="{
                  'message__content--image-and-text':
                    msg.attachment && msg.content,
                  'message__content--only-image':
                    msg.attachment && !msg.content,
                  'message__content--margin-top':
                    conversation.type === 'group' &&
                    msg.sender.userId !== user.userId &&
                    (msg.attachment && !msg.content || msg.attachment && msg.content)
                }"
              >
                <template v-if="msg.attachment && !msg.content">
                  <div class="message__image-padding">
                    <img
                      :src="backendBaseUrl + msg.attachment"
                      class="message__attachment"
                      alt="Attachment"
                    >
                    <span class="text-caption message__image-timestamp">
                      <span
                        v-if="msg.isForwarded"
                        style="font-style: italic"
                      >(forwarded)</span>
                      <span class="text-caption__time-and-ticks">
                        {{
                          formatTime(
                            isEdited(msg)
                              ? msg.editedAt
                              : msg.sentAt,
                          )
                        }}
                        <span
                          v-if="msg.sender.userId === user.userId"
                          class="message__ticks"
                        >
                          <template v-if="isMessageRead(msg)">
                            <svg
                              viewBox="0 0 560 310"
                              fill="none"
                              class="double-tick-icon"
                            >
                              <path
                                d="M40.9705 141.03C31.598 131.658 16.4018 131.658 7.02936 141.03C-2.34312 150.403 -2.34312 165.597 7.02936 174.97L40.9705 141.03ZM152 286L135.03 302.97C144.403 312.342 159.597 312.342 168.97 302.97L152 286ZM424.97 46.9706C434.342 37.5981 434.342 22.402 424.97 13.0295C415.597 3.657 400.403 3.657 391.03 13.0295L424.97 46.9706ZM7.02936 174.97L135.03 302.97L168.97 269.03L40.9705 141.03L7.02936 174.97ZM168.97 302.97L424.97 46.9706L391.03 13.0295L135.03 269.03L168.97 302.97Z"
                                fill="var(--color-tertiary)"
                              />
                              <path
                                d="M168.97 135.03C159.598 125.658 144.402 125.658 135.029 135.03C125.657 144.403 125.657 159.597 135.029 168.97L168.97 135.03ZM280 280L263.03 296.97C272.403 306.342 287.597 306.342 296.97 296.97L280 280ZM552.97 40.9705C562.342 31.598 562.342 16.4018 552.97 7.02936C543.597 -2.34312 528.403 -2.34312 519.03 7.02936L552.97 40.9705ZM135.029 168.97L263.03 296.97L296.97 263.03L168.97 135.03L135.029 168.97ZM296.97 296.97L552.97 40.9705L519.03 7.02936L263.03 263.03L296.97 296.97Z"
                                fill="var(--color-tertiary)"
                              />
                            </svg>
                          </template>
                          <template v-else-if="isMessageSent(msg)">
                            <svg
                              viewBox="0 0 432 304"
                              fill="none"
                              class="single-tick-icon"
                            >
                              <path
                                d="M40.9705 135.03C31.5981 125.658 16.4019 125.658 7.02942 135.03C-2.34306 144.403 -2.34306 159.597 7.02942 168.97L40.9705 135.03ZM152 280L135.03 296.97C144.403 306.342 159.597 306.342 168.97 296.97L152 280ZM424.97 40.9706C434.342 31.5981 434.342 16.402 424.97 7.02948C415.597 -2.343 400.403 -2.343 391.03 7.02948L424.97 40.9706ZM7.02942 168.97L135.03 296.97L168.97 263.03L40.9705 135.03L7.02942 168.97ZM168.97 296.97L424.97 40.9706L391.03 7.02948L135.03 263.03L168.97 296.97Z"
                                fill="var(--color-tertiary)"
                              />
                            </svg>
                          </template>
                        </span>
                      </span>
                    </span>
                  </div>
                </template>
                <template v-else-if="msg.attachment && msg.content">
                  <div class="message__image-padding">
                    <img
                      :ref="(el) => setImageRef(msg.messageId, el)"
                      :src="backendBaseUrl + msg.attachment"
                      class="message__attachment"
                      alt="Attachment"
                      @load="syncMessageWidth(msg.messageId)"
                    >
                  </div>
                  <div class="message__text-and-time">
                    <span class="text-body">{{ msg.content }}</span>
                    <span class="text-caption">
                      <span
                        v-if="msg.isForwarded"
                        style="font-style: italic"
                      >(forwarded)</span>
                      <span class="text-caption__time-and-ticks">
                        {{
                          formatTime(
                            isEdited(msg)
                              ? msg.editedAt
                              : msg.sentAt,
                          )
                        }}
                        <span
                          v-if="msg.sender.userId === user.userId"
                          class="message__ticks"
                        >
                          <template v-if="isMessageRead(msg)">
                            <svg
                              viewBox="0 0 560 310"
                              fill="none"
                              class="double-tick-icon"
                            >
                              <path
                                d="M40.9705 141.03C31.598 131.658 16.4018 131.658 7.02936 141.03C-2.34312 150.403 -2.34312 165.597 7.02936 174.97L40.9705 141.03ZM152 286L135.03 302.97C144.403 312.342 159.597 312.342 168.97 302.97L152 286ZM424.97 46.9706C434.342 37.5981 434.342 22.402 424.97 13.0295C415.597 3.657 400.403 3.657 391.03 13.0295L424.97 46.9706ZM7.02936 174.97L135.03 302.97L168.97 269.03L40.9705 141.03L7.02936 174.97ZM168.97 302.97L424.97 46.9706L391.03 13.0295L135.03 269.03L168.97 302.97Z"
                                fill="var(--color-tertiary)"
                              />
                              <path
                                d="M168.97 135.03C159.598 125.658 144.402 125.658 135.029 135.03C125.657 144.403 125.657 159.597 135.029 168.97L168.97 135.03ZM280 280L263.03 296.97C272.403 306.342 287.597 306.342 296.97 296.97L280 280ZM552.97 40.9705C562.342 31.598 562.342 16.4018 552.97 7.02936C543.597 -2.34312 528.403 -2.34312 519.03 7.02936L552.97 40.9705ZM135.029 168.97L263.03 296.97L296.97 263.03L168.97 135.03L135.029 168.97ZM296.97 296.97L552.97 40.9705L519.03 7.02936L263.03 263.03L296.97 296.97Z"
                                fill="var(--color-tertiary)"
                              />
                            </svg>
                          </template>
                          <template v-else-if="isMessageSent(msg)">
                            <svg
                              viewBox="0 0 432 304"
                              fill="none"
                              class="single-tick-icon"
                            >
                              <path
                                d="M40.9705 135.03C31.5981 125.658 16.4019 125.658 7.02942 135.03C-2.34306 144.403 -2.34306 159.597 7.02942 168.97L40.9705 135.03ZM152 280L135.03 296.97C144.403 306.342 159.597 306.342 168.97 296.97L152 280ZM424.97 40.9706C434.342 31.5981 434.342 16.402 424.97 7.02948C415.597 -2.343 400.403 -2.343 391.03 7.02948L424.97 40.9706ZM7.02942 168.97L135.03 296.97L168.97 263.03L40.9705 135.03L7.02942 168.97ZM168.97 296.97L424.97 40.9706L391.03 7.02948L135.03 263.03L168.97 296.97Z"
                                fill="var(--color-tertiary)"
                              />
                            </svg>
                          </template>
                        </span>
                      </span>
                      <span v-if="isEdited(msg)" style="font-style: italic">(edited)</span>
                    </span>
                  </div>
                </template>
                <template v-else>
                  <div class="message__only-text">
                    <span class="text-body">{{ msg.content }}</span>
                    <span class="text-caption">
                      <span
                        v-if="msg.isForwarded"
                        style="font-style: italic"
                      >(forwarded)</span>
                      <span class="text-caption__time-and-ticks">
                        {{
                          formatTime(
                            isEdited(msg)
                              ? msg.editedAt
                              : msg.sentAt,
                          )
                        }}
                        <span
                          v-if="msg.sender.userId === user.userId"
                          class="message__ticks"
                        >
                          <template v-if="isMessageRead(msg)">
                            <svg
                              viewBox="0 0 560 310"
                              fill="none"
                              class="double-tick-icon"
                            >
                              <path
                                d="M40.9705 141.03C31.598 131.658 16.4018 131.658 7.02936 141.03C-2.34312 150.403 -2.34312 165.597 7.02936 174.97L40.9705 141.03ZM152 286L135.03 302.97C144.403 312.342 159.597 312.342 168.97 302.97L152 286ZM424.97 46.9706C434.342 37.5981 434.342 22.402 424.97 13.0295C415.597 3.657 400.403 3.657 391.03 13.0295L424.97 46.9706ZM7.02936 174.97L135.03 302.97L168.97 269.03L40.9705 141.03L7.02936 174.97ZM168.97 302.97L424.97 46.9706L391.03 13.0295L135.03 269.03L168.97 302.97Z"
                                fill="var(--color-tertiary)"
                              />
                              <path
                                d="M168.97 135.03C159.598 125.658 144.402 125.658 135.029 135.03C125.657 144.403 125.657 159.597 135.029 168.97L168.97 135.03ZM280 280L263.03 296.97C272.403 306.342 287.597 306.342 296.97 296.97L280 280ZM552.97 40.9705C562.342 31.598 562.342 16.4018 552.97 7.02936C543.597 -2.34312 528.403 -2.34312 519.03 7.02936L552.97 40.9705ZM135.029 168.97L263.03 296.97L296.97 263.03L168.97 135.03L135.029 168.97ZM296.97 296.97L552.97 40.9705L519.03 7.02936L263.03 263.03L296.97 296.97Z"
                                fill="var(--color-tertiary)"
                              />
                            </svg>
                          </template>
                          <template v-else-if="isMessageSent(msg)">
                            <svg
                              viewBox="0 0 432 304"
                              fill="none"
                              class="single-tick-icon"
                            >
                              <path
                                d="M40.9705 135.03C31.5981 125.658 16.4019 125.658 7.02942 135.03C-2.34306 144.403 -2.34306 159.597 7.02942 168.97L40.9705 135.03ZM152 280L135.03 296.97C144.403 306.342 159.597 306.342 168.97 296.97L152 280ZM424.97 40.9706C434.342 31.5981 434.342 16.402 424.97 7.02948C415.597 -2.343 400.403 -2.343 391.03 7.02948L424.97 40.9706ZM7.02942 168.97L135.03 296.97L168.97 263.03L40.9705 135.03L7.02942 168.97ZM168.97 296.97L424.97 40.9706L391.03 7.02948L135.03 263.03L168.97 296.97Z"
                                fill="var(--color-tertiary)"
                              />
                            </svg>
                          </template>
                        </span>
                      </span>
                      <span v-if="isEdited(msg)" style="font-style: italic">(edited)</span>
                    </span>
                  </div>
                </template>
                <button
                  v-if="hoveredMessageId === msg.messageId"
                  :ref="(el) => setButtonRef(msg.messageId, el)"
                  class="message__dropdown-menu-button"
                  @click.stop="openMenu(msg.messageId, $event)"
                >
                  <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    class="message__dropdown-menu-icon"
                  >
                    <path
                      d="M7 10L12 15L17 10"
                      stroke="var(--color-tertiary)"
                      stroke-width="1.5"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                </button>
                <div
                  v-if="menuOpenFor === msg.messageId"
                  ref="menuRef"
                  class="message__dropdown-menu"
                  :style="menuStyles"
                >
                  <button @click="commentMessage(msg)">
                    <svg viewBox="0 0 24 24" fill="none">
                      <path
                        d="M8.5 11C9.32843 11 10 10.3284 10 9.5C10 8.67157 9.32843 8 8.5 8C7.67157 8 7 8.67157 7 9.5C7 10.3284 7.67157 11 8.5 11Z"
                        fill="var(--color-tertiary)"
                      />
                      <path
                        d="M17 9.5C17 10.3284 16.3284 11 15.5 11C14.6716 11 14 10.3284 14 9.5C14 8.67157 14.6716 8 15.5 8C16.3284 8 17 8.67157 17 9.5Z"
                        fill="var(--color-tertiary)"
                      />
                      <path
                        d="M8.88875 13.5414C8.63822 13.0559 8.0431 12.8607 7.55301 13.1058C7.05903 13.3528 6.8588 13.9535 7.10579 14.4474C7.18825 14.6118 7.29326 14.7659 7.40334 14.9127C7.58615 15.1565 7.8621 15.4704 8.25052 15.7811C9.04005 16.4127 10.2573 17.0002 12.0002 17.0002C13.7431 17.0002 14.9604 16.4127 15.7499 15.7811C16.1383 15.4704 16.4143 15.1565 16.5971 14.9127C16.7076 14.7654 16.8081 14.6113 16.8941 14.4485C17.1387 13.961 16.9352 13.3497 16.4474 13.1058C15.9573 12.8607 15.3622 13.0559 15.1117 13.5414C15.0979 13.5663 14.9097 13.892 14.5005 14.2194C14.0401 14.5877 13.2573 15.0002 12.0002 15.0002C10.7431 15.0002 9.96038 14.5877 9.49991 14.2194C9.09071 13.892 8.90255 13.5663 8.88875 13.5414Z"
                        fill="var(--color-tertiary)"
                      />
                      <path
                        fill-rule="evenodd"
                        clip-rule="evenodd"
                        d="M12 23C18.0751 23 23 18.0751 23 12C23 5.92487 18.0751 1 12 1C5.92487 1 1 5.92487 1 12C1 18.0751 5.92487 23 12 23ZM12 20.9932C7.03321 20.9932 3.00683 16.9668 3.00683 12C3.00683 7.03321 7.03321 3.00683 12 3.00683C16.9668 3.00683 20.9932 7.03321 20.9932 12C20.9932 16.9668 16.9668 20.9932 12 20.9932Z"
                        fill="var(--color-tertiary)"
                      />
                    </svg>
                    <span class="text-body">Comment</span>
                  </button>
                  <button @click="forwardMessage(msg)">
                    <svg viewBox="0 0 24 24" fill="none">
                      <path
                        d="M4 17V15.8C4 14.1198 4 13.2798 4.32698 12.638C4.6146 12.0735 5.07354 11.6146 5.63803 11.327C6.27976 11 7.11984 11 8.8 11H20M20 11L16 7M20 11L16 15"
                        stroke="var(--color-tertiary)"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                    <span class="text-body">Forward</span>
                  </button>
                  <button
                    v-if="
                      msg.sender.userId === user.userId &&
                        !(msg.attachment && !msg.content) &
                        !msg.isForwarded
                    "
                    @click="editMessage(msg)"
                  >
                    <svg viewBox="0 0 24 24" fill="none">
                      <path
                        d="M20.1498 7.93997L8.27978 19.81C7.21978 20.88 4.04977 21.3699 3.32977 20.6599C2.60977 19.9499 3.11978 16.78 4.17978 15.71L16.0498 3.84C16.5979 3.31801 17.3283 3.03097 18.0851 3.04019C18.842 3.04942 19.5652 3.35418 20.1004 3.88938C20.6356 4.42457 20.9403 5.14781 20.9496 5.90463C20.9588 6.66146 20.6718 7.39189 20.1498 7.93997V7.93997Z"
                        stroke="var(--color-tertiary)"
                        stroke-width="1.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                    <span class="text-body">Edit</span>
                  </button>
                  <button
                    v-if="msg.sender.userId === user.userId"
                    @click="deleteMessage(msg)"
                  >
                    <svg viewBox="0 0 24 24" fill="none">
                      <path
                        d="M3 6.98996C8.81444 4.87965 15.1856 4.87965 21 6.98996"
                        stroke="var(--color-tertiary)"
                        stroke-width="1.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M8.00977 5.71997C8.00977 4.6591 8.43119 3.64175 9.18134 2.8916C9.93148 2.14146 10.9489 1.71997 12.0098 1.71997C13.0706 1.71997 14.0881 2.14146 14.8382 2.8916C15.5883 3.64175 16.0098 4.6591 16.0098 5.71997"
                        stroke="var(--color-tertiary)"
                        stroke-width="1.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M12 13V18"
                        stroke="var(--color-tertiary)"
                        stroke-width="1.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M19 9.98999L18.33 17.99C18.2225 19.071 17.7225 20.0751 16.9246 20.8123C16.1266 21.5494 15.0861 21.9684 14 21.99H10C8.91389 21.9684 7.87336 21.5494 7.07541 20.8123C6.27745 20.0751 5.77745 19.071 5.67001 17.99L5 9.98999"
                        stroke="var(--color-tertiary)"
                        stroke-width="1.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                    <span class="text-body">Delete</span>
                  </button>
                </div>
              </div>
            </div>
            <div
              v-if="msg.comments && msg.comments.length"
              class="comment-list"
            >
              <template
                v-for="group in groupComments(msg.comments)"
                :key="group.emoji"
              >
                <button
                  class="comment"
                  :class="{
                    'my-comment': group.users.includes(user.userId),
                  }"
                  @click="onEmojiClick(msg.messageId, group)"
                >
                  {{ group.emoji }}
                  <span v-if="group.users.length > 1" class="text-secondary">{{
                    group.users.length
                  }}</span>
                </button>
              </template>
            </div>
          </div>
        </template>
      </div>
      <div
        v-if="emojiMenuFor"
        class="emoji-menu"
        :style="{
          position: 'fixed',
          top: emojiMenuPosition.top + 'px',
          left: emojiMenuPosition.left + 'px',
          zIndex: 10000,
        }"
      >
        <button
          v-for="emoji in emojiOptions"
          :key="emoji"
          class="emoji-button"
          @click="addEmojiComment(emojiMenuFor, emoji)"
        >
          {{ emoji }}
        </button>
      </div>
    </div>
    <div class="message-field">
      <button class="attachment__button" @click="onAttachmentClick">
        <svg viewBox="0 0 1920 1920" fill="none" class="attchment__icon">
          <path
            d="M866.332 213v653.332H213v186.666h653.332v653.332h186.666v-653.332h653.332V866.332h-653.332V213z"
            fill-rule="evenodd"
            fill="var(--color-secondary)"
          />
        </svg>
      </button>
      <input
        ref="fileInput"
        type="file"
        style="display: none"
        accept="image/*"
        @change="onImageChange"
      >
      <button
        v-if="imagePreviewUrl"
        class="message__attachment-preview"
        @click="removeAttachment"
      >
        <img :src="imagePreviewUrl" alt="Preview">
      </button>
      <textarea
        v-model="message"
        placeholder="Write a message"
        rows="1"
        class="message-input"
        @keydown.enter="onEnter"
        @input="autoResize"
      />
      <button class="send__button" @click="sendMessage">
        <svg viewBox="0 0 28 28" class="send__icon">
          <g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
            <g fill="var(--color-background)" fill-rule="nonzero">
              <path
                d="M3.78963301,2.77233335 L24.8609339,12.8499121 C25.4837277,13.1477699 25.7471402,13.8941055 25.4492823,14.5168992 C25.326107,14.7744476 25.1184823,14.9820723 24.8609339,15.1052476 L3.78963301,25.1828263 C3.16683929,25.4806842 2.42050372,25.2172716 2.12264586,24.5944779 C1.99321184,24.3238431 1.96542524,24.015685 2.04435886,23.7262618 L4.15190935,15.9983421 C4.204709,15.8047375 4.36814355,15.6614577 4.56699265,15.634447 L14.7775879,14.2474874 C14.8655834,14.2349166 14.938494,14.177091 14.9721837,14.0981464 L14.9897199,14.0353553 C15.0064567,13.9181981 14.9390703,13.8084248 14.8334007,13.7671556 L14.7775879,13.7525126 L4.57894108,12.3655968 C4.38011873,12.3385589 4.21671819,12.1952832 4.16392965,12.0016992 L2.04435886,4.22889788 C1.8627142,3.56286745 2.25538645,2.87569101 2.92141688,2.69404635 C3.21084015,2.61511273 3.51899823,2.64289932 3.78963301,2.77233335 Z"
              />
            </g>
          </g>
        </svg>
      </button>
    </div>
  </div>
  <div v-else class="conversation" />
</template>

<script setup>
import {
  ref,
  nextTick,
  watch,
  onBeforeUnmount,
  computed,
  watchEffect,
} from "vue";
import api from "@/services/api";
import { backendBaseUrl } from "@/services/baseUrl";

import defaultProfilePicture from "@/assets/default-profile-picture.jpg";
import defaultGroupPicture from "@/assets/default-group-picture.jpg";

const props = defineProps({
  user: Object,
  conversation: Object,
  forwardModalOpen: Boolean,
  sidebarOpen: Boolean,
});

const emit = defineEmits([
  "message-sent",
  "conversation-updated",
  "show-info-sidebar",
  "forward-modal-open",
]);

function getOtherUser(conversation) {
  if (conversation.type !== "private") return null;
  return conversation.participants.find(
    (user) => user.userId !== props.user?.userId,
  );
}

const message = ref("");
const imageFile = ref(null);
const fileInput = ref(null);

const messageRefs = ref({});
const imageRefs = ref({});

function setMessageRef(messageId, el) {
  if (el) messageRefs.value[messageId] = el;
}
function setImageRef(messageId, el) {
  if (el) imageRefs.value[messageId] = el;
}

function syncMessageWidth(messageId) {
  nextTick(() => {
    const img = imageRefs.value[messageId];
    const msg = messageRefs.value[messageId];

    if (editingMessageId.value === messageId) {
      if (msg) msg.style.maxWidth = "";
      return;
    }

    if (img && msg) {
      msg.style.maxWidth = "";
      let parent = msg.parentElement;
      while (parent && !parent.classList.contains("messages-wrapper")) {
        parent = parent.parentElement;
      }
      const parentWidth = parent ? parent.offsetWidth : window.innerWidth;
      const maxWidth = parentWidth * 0.6;

      if (img.clientWidth < maxWidth) {
        msg.style.maxWidth = img.clientWidth + "px";
      } else {
        msg.style.maxWidth = "60%";
      }
    }
  });
}

function onEnter(e) {
  if (!e.shiftKey) {
    e.preventDefault();
    sendMessage();
  }
}

function onAttachmentClick() {
  if (fileInput.value) {
    fileInput.value.click();
  }
}

const imagePreviewUrl = ref(null);

function onImageChange(e) {
  imageFile.value = e.target.files[0];
  if (imageFile.value) {
    imagePreviewUrl.value = URL.createObjectURL(imageFile.value);
  } else {
    imagePreviewUrl.value = null;
  }
}

function removeAttachment() {
  imageFile.value = null;
  imagePreviewUrl.value = null;
  if (fileInput.value) fileInput.value.value = "";
}

async function sendMessage() {
  if (message.value.trim() === "" && !imageFile.value) return;

  const formData = new FormData();
  if (message.value.trim()) formData.append("content", message.value.trim());
  if (imageFile.value) formData.append("image", imageFile.value);

  try {
    const response = await api.post(
      `/conversations/${props.conversation.conversationId}/messages`,
      formData,
      { headers: { "Content-Type": "multipart/form-data" } },
    );

    if (!Array.isArray(messages.value)) {
      messages.value = [];
    }

    messages.value.push(response.data);
    message.value = "";
    imageFile.value = null;
    imagePreviewUrl.value = null;

    nextTick(autoResize);

    emit("message-sent");
  } catch (e) {
    console.error(e);
  }

  message.value = "";
  nextTick(autoResize);
}

function autoResize(e) {
  const el = e ? e.target : document.querySelector(".message-input");
  if (el) {
    el.style.height = "auto";
    el.style.height = el.scrollHeight + "px";
  }
}

function autoResizeEdit(e) {
  const el = e ? e.target : editInput.value;
  if (el) {
    el.style.height = "auto";
    el.style.height = el.scrollHeight + "px";
  }
}

const messages = ref([]);

const messagesContainer = ref(null);

watch(
  messages,
  async () => {
    await nextTick();
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
    }
  },
  { deep: true },
);

function formatTime(sentAt) {
  if (!sentAt) return "";

  const date = new Date(sentAt);
  return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
}

function formatDay(dateString) {
  const date = new Date(dateString);
  const now = new Date();

  if (
    date.getFullYear() === now.getFullYear() &&
    date.getMonth() === now.getMonth() &&
    date.getDate() === now.getDate()
  ) {
    return "Today";
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

function isSameDay(date1, date2) {
  if (!date1 || !date2) return false;
  const d1 = new Date(date1);
  const d2 = new Date(date2);
  return (
    d1.getFullYear() === d2.getFullYear() &&
    d1.getMonth() === d2.getMonth() &&
    d1.getDate() === d2.getDate()
  );
}

async function fetchMessages(conversationId) {
  if (!conversationId) {
    messages.value = [];
    return;
  }

  try {
    const response = await api.get(`/conversations/${conversationId}`);
    messages.value = response.data.messages;
  } catch (e) {
    messages.value = [];
    console.error(e);
  }
}

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
  }
}

function waitForImagesToLoad() {
  const imgs = messagesContainer.value?.querySelectorAll(
    ".message__attachment",
  );
  if (!imgs || imgs.length === 0) {
    scrollToBottom();
    return;
  }
  let loaded = 0;
  imgs.forEach((img) => {
    if (img.complete) {
      loaded++;
    } else {
      img.addEventListener(
        "load",
        () => {
          loaded++;
          if (loaded === imgs.length) {
            scrollToBottom();
          }
        },
        { once: true },
      );
    }
  });
  if (loaded === imgs.length) {
    scrollToBottom();
  }
}

watch(
  () => props.conversation?.conversationId,
  async (newId) => {
    if (newId && props.conversation) {
      await fetchMessages(newId);
      await nextTick();
      waitForImagesToLoad();
    }
  },
  { immediate: true },
);

function isEdited(message) {
  return !!message.editedAt && message.editedAt !== "0001-01-01T00:00:00Z";
}

function isMessageRead(message) {
  if (!message.trackings || !message.trackings.read) return false;

  const others =
    props.conversation.type === "group"
      ? props.conversation.members.filter(
          (u) => u.userId !== message.sender.userId,
        )
      : props.conversation.participants.filter(
          (u) => u.userId !== message.sender.userId,
        );

  return others.every((user) => {
    const readAt = message.trackings.read[user.userId];
    return !!readAt;
  });
}

function isMessageSent(message) {
  return !!message.messageId;
}

const hoveredMessageId = ref(null);
const menuOpenFor = ref(null);
const menuRef = ref(null);

const buttonRefs = ref({});
const menuPosition = ref({ top: 0, left: 0 });

function setButtonRef(messageId, el) {
  if (el) buttonRefs.value[messageId] = el;
}

function openMenu(messageId) {
  menuOpenFor.value = messageId;
  nextTick(() => {
    const btn = buttonRefs.value[messageId];
    if (btn) {
      const rect = btn.getBoundingClientRect();
      const menuWidth = 180;
      const menuHeight = 160;
      const margin = 32;
      let top = rect.bottom + 4;
      let left = rect.right - menuWidth;

      if (top + menuHeight > window.innerHeight - margin) {
        top = rect.top - menuHeight - 4;
      }

      if (left < margin) {
        left = rect.left;
      }

      if (left + menuWidth > window.innerWidth - margin) {
        left = window.innerWidth - menuWidth - 8;
      }

      menuPosition.value = { top, left };
    }

    document.addEventListener("mousedown", handleClickOutside);
  });
}

const menuStyles = computed(() => ({
  top: menuPosition.value.top + "px",
  left: menuPosition.value.left + "px",
}));

function closeMenu() {
  menuOpenFor.value = null;
  document.removeEventListener("mousedown", handleClickOutside);
}

function handleClickOutside(event) {
  if (Array.isArray(menuRef.value)) {
    const clickedInside = menuRef.value.some(
      (el) => el && el.contains(event.target),
    );
    if (!clickedInside) closeMenu();
  } else if (menuRef.value && !menuRef.value.contains(event.target)) {
    closeMenu();
  }
}

const editingMessageId = ref(null);
const editingContent = ref("");
const editInput = ref(null);

async function editMessage(message) {
  if (message.attachment && !message.content) return;

  closeMenu();
  editingMessageId.value = message.messageId;
  editingContent.value = message.content;
}

watch(editingMessageId, async (newVal, oldVal) => {
  if (newVal) {
    await nextTick();

    const inputEl = Array.isArray(editInput.value)
      ? editInput.value[0]
      : editInput.value;
    if (inputEl && typeof inputEl.focus === "function") {
      inputEl.focus();
      autoResizeEdit({ target: inputEl });
    }
  }

  if (newVal && messageRefs.value[newVal]) {
    messageRefs.value[newVal].style.maxWidth = "";
  }

  if (oldVal && imageRefs.value[oldVal]) {
    syncMessageWidth(oldVal);
  }
});

function onEditEnter(e, message) {
  if (!e.shiftKey) {
    e.preventDefault();
    saveEdit(message);
  }
}

async function saveEdit(message) {
  if ((message.content ?? "") === (editingContent.value ?? "")) {
    editingMessageId.value = null;
    editingContent.value = "";
    return;
  }

  try {
    const response = await api.put(`/messages/${message.messageId}`, {
      content: editingContent.value,
    });

    const idx = messages.value.findIndex(
      (m) => m.messageId === message.messageId,
    );

    if (idx !== -1) messages.value[idx] = response.data;

    editingMessageId.value = null;
    editingContent.value = "";
  } catch (e) {
    console.error(e);
  }
}

function cancelEdit() {
  editingMessageId.value = null;
  editingContent.value = "";
}

async function deleteMessage(message) {
  try {
    await api.delete(`/messages/${message.messageId}`);
    messages.value = messages.value.filter(
      (m) => m.messageId !== message.messageId,
    );

    closeMenu();
  } catch (e) {
    console.error(e);
  }
}

function forwardMessage(message) {
  emit("forward-modal-open", { open: true, message });
  closeMenu();
}

const emojiMenuFor = ref(null);
const emojiOptions = ["👍", "😂", "❤️", "😮"];
const pendingEmojiMenuPosition = ref(null);
const emojiMenuPosition = ref({ top: 0, left: 0 });

function commentMessage(message, event) {
  let menuEl = menuRef.value;
  if (Array.isArray(menuEl)) menuEl = menuEl[0];
  let position = { top: 100, left: 100 };

  if (menuEl) {
    const rect = menuEl.getBoundingClientRect();
    position = {
      top: rect.bottom + 4,
      left: rect.left,
    };
  }

  pendingEmojiMenuPosition.value = position;
  closeMenu();

  nextTick(() => {
    openEmojiMenu(message.messageId);
  });
}

function openEmojiMenu(messageId) {
  emojiMenuFor.value = messageId;

  if (pendingEmojiMenuPosition.value) {
    emojiMenuPosition.value = pendingEmojiMenuPosition.value;
    pendingEmojiMenuPosition.value = null;
  } else {
    emojiMenuPosition.value = { top: 100, left: 100 };
  }

  document.addEventListener("mousedown", handleClickOutsideEmojiMenu);
}

function closeEmojiMenu() {
  emojiMenuFor.value = null;
  document.removeEventListener("mousedown", handleClickOutsideEmojiMenu);
}

function handleClickOutsideEmojiMenu(event) {
  if (!event.target.closest(".emoji-menu")) {
    closeEmojiMenu();
  }
}

async function addEmojiComment(messageId, emoji) {
  const message = messages.value.find((m) => m.messageId === messageId);
  if (!message) return;

  const alreadyReacted = (message.comments || []).some(
    (c) => c.commenter.userId === props.user.userId,
  );

  if (alreadyReacted) return;

  await api.post(`/messages/${messageId}/comments`, { emoji: emoji });

  await fetchMessages(props.conversation.conversationId);
  closeEmojiMenu();
}

watch(
  () => menuOpenFor.value !== null || emojiMenuFor.value !== null,
  (isMenuOpen) => {
    const chatWrapper = document.querySelector(".messages-wrapper");
    if (chatWrapper) {
      if (isMenuOpen) {
        chatWrapper.classList.add("no-scroll");
      } else {
        chatWrapper.classList.remove("no-scroll");
      }
    }
  },
);

onBeforeUnmount(() => {
  document.removeEventListener("mousedown", handleClickOutside);
});

function groupComments(comments) {
  const grouped = {};

  for (const comment of comments) {
    const userId = comment.commenter?.userId;
    if (!grouped[comment.emoji]) {
      grouped[comment.emoji] = {
        emoji: comment.emoji,
        users: [],
        comments: [],
      };
    }
    if (userId) {
      grouped[comment.emoji].users.push(userId);
    }
    grouped[comment.emoji].comments.push(comment);
  }

  return Object.values(grouped);
}

async function uncommentMessage(comment) {
  try {
    await api.delete(`/comments/${comment.commentId}`);
    await fetchMessages(props.conversation.conversationId);
  } catch (e) {
    console.error(e);
  }
}

function onEmojiClick(messageId, group) {
  if (group.users.includes(props.user.userId)) {
    const myComment = group.comments.find(
      (c) => c.commenter.userId === props.user.userId,
    );

    if (myComment) {
      uncommentMessage(myComment);
    }
  } else {
    addEmojiComment(messageId, group.emoji);
  }
}
</script>

<style scoped>
.conversation {
  display: flex;
  flex-direction: column;
  width: 70vw;
  height: 100vh;
}

.conversation--sidebar-open {
  width: 40vw;
}

.conversation-header {
  width: 100%;
}

.conversation-header__button {
  display: flex;
  align-items: center;
  gap: 1rem;
  width: 100%;
  padding: 1rem;
  border: none;
  background-color: inherit;
  color: var(--color-secondary);
}

.conversation-header__photo {
  width: 40px;
  height: 40px;
  border: 2px transparent;
  border-radius: 50%;
  object-fit: cover;
}

.messages-wrapper {
  flex: 1 1 auto;
  overflow-y: auto;
  min-height: 0;
}

.messages {
  position: relative;
  display: flex;
  justify-content: flex-end;
  flex-direction: column;
  gap: 0.25rem;
  min-height: 100%;
  padding: 0rem 3rem;
}

.date-separator {
  display: flex;
  justify-content: center;
  align-items: center;
  align-self: center;
  padding: 0.25rem 0.5rem;
  margin: 0.5rem 0rem;
  border-radius: 8px;
  background-color: var(--color-quaternary);
}

.message-wrapper {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.25rem;
}

.message-wrapper--mine {
  align-items: flex-end;
}

.message__avatar {
  width: 24px;
  height: 24px;
  border: 1px solid var(--color-quaternary);
  border-radius: 50%;
  object-fit: cover;
}

.message {
  position: relative;
  display: flex;
  flex-direction: column;
  border-radius: 8px;
  padding: 0.25rem;
  background-color: var(--color-quaternary);
  max-width: 60%;
  width: fit-content;
}

.message--mine {
  background-color: var(--color-primary);
  align-items: flex-end;
}

.message--mine.editing {
  max-width: 100%;
  flex: 1;
  background: none;
  border: 3px solid var(--color-primary);
}

.sender-name {
  padding: 0rem 0.25rem;
}

.message__content {
  display: flex;
  align-items: flex-end;
  gap: 0.5rem;
  white-space: pre-line;
  overflow-wrap: anywhere;
  background-color: inherit;
}

.message__content--image-and-text {
  flex-direction: column;
  align-items: flex-start;
  gap: 0.5rem;
}

.message__image-padding {
  position: relative;
  width: auto;
}

.message__image-timestamp {
  position: absolute;
  right: 8px;
  bottom: 8px;
  background: var(--color-quaternary);
  color: var(--color-tertiary);
  padding: 0.25rem 0.5rem;
  max-width: 80px;
  border-radius: 8px;
  font-size: 0.7rem;
  pointer-events: none;
}

.message__only-text {
  display: flex;
  align-items: flex-end;
  gap: 0.5rem;
  white-space: pre-line;
  overflow-wrap: anywhere;
  background-color: inherit;
  padding: 0.25rem;
}

.message__text-and-time {
  display: flex;
  justify-content: space-between;
  width: 100%;
  align-items: flex-end;
  gap: 0.5rem;
  padding: 0.25rem;
}

.message__text-and-time .text-body {
  flex: 1 1 auto;
  text-align: left;
}

.message__text-and-time .text-caption {
  text-align: right;
}

.text-caption__time-and-ticks {
  display: flex;
  gap: 0.25rem;
}

.message--mine.editing .message__content {
  width: 100%;
  flex: 1 1 auto;
}

.message__content .text-body {
  white-space: pre-line;
  word-break: break-word;
  overflow-wrap: anywhere;
  flex: 1 1 auto;
  min-width: 0;
}

.message__attachment {
  width: 100%;
  max-height: 380px;
  border-radius: 8px;
  object-fit: cover;
  display: block;
  width: 100%;
}

.message__content .text-caption {
  display: flex;
  align-items: flex-end;
  justify-content: flex-end;
  gap: 0.25rem;
  white-space: normal;
  flex-shrink: 1;
  flex-wrap: wrap;
  word-break: break-word;
  overflow-wrap: anywhere;
  min-width: 70px;
}

.message__edit {
  background-color: inherit;
  font-size: 1rem;
  color: var(--color-secondary);
  resize: none;
  outline: none;
  line-height: 1.5;
  font-family: inherit;
  border: none;
}

.message__ticks {
  display: flex;
  justify-content: flex-end;
}

.single-tick-icon {
  width: 11px;
  height: 11px;
}

.double-tick-icon {
  width: 14px;
  height: 11px;
}

.message__dropdown-menu-button {
  position: absolute;
  top: 4px;
  right: 0px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: inherit;
  border: none;
  border-radius: 50%;
  z-index: 2;
  box-shadow: -1px 1px 4px 0 var(--color-quaternary);
}

.message__content--margin-top {
  margin-top: 0.5rem;
}

.message--mine > .message__content--margin-top,
.conversation[type="private"] .message__content--margin-top {
  margin-top: 0rem;
}

.message__content--image-and-text .message__dropdown-menu-button,
.message__content--only-image .message__dropdown-menu-button {
  right: 4px;
  background: none;
  box-shadow: none !important;
}

.message--mine .message__dropdown-menu-button {
  box-shadow: -1px 1px 4px 0 var(--color-primary);
}

.message__dropdown-menu-icon {
  width: 24px;
  height: 24px;
}

.message__dropdown-menu {
  position: fixed;
  background-color: #252526;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  z-index: 9999;
  padding: 0.5rem;
}

.message__dropdown-menu button {
  display: flex;
  align-items: center;
  gap: 1rem;
  background-color: inherit;
  border: none;
  border-radius: 8px;
  padding: 0.5rem;
  transition: background 0.1s;
}

.message__dropdown-menu button:hover {
  background-color: var(--color-quaternary);
}

.message__dropdown-menu button svg {
  width: 24px;
  height: 24px;
}

.message__dropdown-menu button .text-body {
  color: var(--color-tertiary);
}

.comment-list {
  display: flex;
  gap: 0.25rem;
}

.comment {
  border: none;
  background-color: var(--color-quaternary);
  border-radius: 8px;
  padding: 0.25rem 0.5rem;
  gap: 0.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.1s;
}

.comment.my-comment {
  background-color: var(--color-primary);
}

.emoji-menu {
  background-color: #252526;
  border-radius: 50px;
  display: flex;
  padding: 0.5rem;
}

.emoji-button {
  display: flex;
  align-items: center;
  gap: 1rem;
  background-color: inherit;
  border: none;
  border-radius: 50%;
  padding: 0rem 0.4rem;
  font-size: 1.75rem;
}

.message-field {
  display: flex;
  align-items: flex-end;
  gap: 0.25rem;
  margin: 1rem;
  padding: 0.25rem;
  padding-left: 0.5rem;
  border-radius: 24px;
  background-color: var(--color-quaternary);
}

.attachment__button {
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 24px;
  margin-bottom: 3.2px;
  padding: 0.3rem;
  background-color: inherit;
}

.attachment__button:hover {
  background-color: #252526;
  transition: background 0.1s;
}

.attchment__icon {
  width: 20px;
  height: 20px;
}

.message__attachment-preview {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 8px;
  overflow: hidden;
  background-color: inherit;
}

.message__attachment-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
}

.message-input {
  border: none;
  background-color: inherit;
  font-size: 1rem;
  color: var(--color-secondary);
  flex: 1;
  min-width: 0;
  resize: none;
  outline: none;
  margin-bottom: 6px;
  padding: 0;
  font-family: inherit;
  line-height: 1.5;
}

.send__button {
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 24px;
  padding: 0.5rem;
  background-color: var(--color-primary);
  transition: filter 0.1s;
}

.send__button:hover {
  filter: brightness(1.05);
}

.send__icon {
  width: 20px;
  height: 20px;
}

.messages-wrapper.no-scroll {
  overflow: hidden !important;
}
</style>
