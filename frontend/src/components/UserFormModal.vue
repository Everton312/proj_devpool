<script setup>
import { ref, watch } from 'vue'
import { userService } from '../services/api'

const props = defineProps({
  userToEdit: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['saved', 'close'])

const name = ref('')
const email = ref('')
const login = ref('')
const password = ref('')
const currentPassword = ref('')
const newPassword = ref('')
const active = ref(true)

const loading = ref(false)
const errorMessage = ref('')

watch(
  () => props.userToEdit,
  (user) => {
    errorMessage.value = ''
    if (user) {
      name.value = user.name
      email.value = user.email
      login.value = user.login
      active.value = user.active
      currentPassword.value = ''
      newPassword.value = ''
    } else {
      name.value = ''
      email.value = ''
      login.value = ''
      password.value = ''
      active.value = true
    }
  },
  { immediate: true }
)

async function handleSubmit() {
  loading.value = true
  errorMessage.value = ''

  try {
    if (props.userToEdit) {
      const payload = {
        name: name.value,
        email: email.value,
        login: login.value,
        active: active.value
      }
      if (newPassword.value) {
        payload.current_password = currentPassword.value
        payload.new_password = newPassword.value
      }
      await userService.update(props.userToEdit.id, payload)
      emit('saved', 'Usuário atualizado com sucesso!')
    } else {
      await userService.create({
        name: name.value,
        email: email.value,
        login: login.value,
        password: password.value,
        active: active.value
      })
      emit('saved', 'Usuário cadastrado com sucesso!')
    }
  } catch (err) {
    errorMessage.value = err.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="card border-0 shadow-sm rounded-4 overflow-hidden mx-auto" style="max-width: 680px;">
    <div class="card-header bg-white border-bottom py-4 px-4 d-flex justify-content-between align-items-center">
      <div>
        <h2 class="h4 fw-bold mb-1 text-dark">
          {{ userToEdit ? 'Editar usuário' : 'Novo usuário' }}
        </h2>
        <p class="text-muted small mb-0">
          {{ userToEdit ? 'Atualize os dados cadastrais abaixo.' : 'Preencha os campos para cadastrar um novo usuário.' }}
        </p>
      </div>

      <button
        type="button"
        class="btn btn-light btn-sm border px-3 py-2 rounded-3"
        :disabled="loading"
        @click="emit('close')"
      >
        ← Voltar
      </button>
    </div>

    <div class="card-body p-4">
      <div v-if="errorMessage" class="alert alert-danger py-2 rounded-3" role="alert">
        {{ errorMessage }}
      </div>

      <form @submit.prevent="handleSubmit" class="row g-3">
        <div class="col-12">
          <label class="form-label fw-medium small text-secondary">Nome completo</label>
          <input
            v-model="name"
            type="text"
            class="form-control rounded-3"
            placeholder="Ex: Maria Silva"
            required
            minlength="3"
          />
        </div>

        <div class="col-md-6">
          <label class="form-label fw-medium small text-secondary">E-mail</label>
          <input
            v-model="email"
            type="email"
            class="form-control rounded-3"
            placeholder="usuario@dominio.com"
            required
          />
        </div>

        <div class="col-md-6">
          <label class="form-label fw-medium small text-secondary">Login</label>
          <input
            v-model="login"
            type="text"
            class="form-control rounded-3"
            placeholder="Mínimo 3 caracteres"
            required
            minlength="3"
          />
        </div>

        <div v-if="!userToEdit" class="col-12">
          <label class="form-label fw-medium small text-secondary">Senha</label>
          <input
            v-model="password"
            type="password"
            class="form-control rounded-3"
            placeholder="Mínimo 8 caracteres"
            required
            minlength="8"
          />
        </div>

        <template v-else>
          <div class="col-12 pt-2">
            <hr class="my-1 text-muted" />
            <span class="text-muted small">Alteração de senha (opcional)</span>
          </div>

          <div class="col-md-6">
            <label class="form-label fw-medium small text-secondary">Senha atual</label>
            <input
              v-model="currentPassword"
              type="password"
              class="form-control rounded-3"
              placeholder="Digite sua senha atual"
            />
          </div>

          <div class="col-md-6">
            <label class="form-label fw-medium small text-secondary">Nova senha</label>
            <input
              v-model="newPassword"
              type="password"
              class="form-control rounded-3"
              placeholder="Mínimo 8 caracteres"
              minlength="8"
            />
          </div>
        </template>

        <div class="col-12 pt-2">
          <div class="form-check form-switch">
            <input
              id="activeSwitch"
              v-model="active"
              class="form-check-input"
              type="checkbox"
              role="switch"
            />
            <label class="form-check-label fw-medium text-dark" for="activeSwitch">
              Usuário ativo no sistema
            </label>
          </div>
        </div>

        <div class="col-12 d-flex justify-content-end gap-2 pt-3 border-top mt-4">
          <button
            type="button"
            class="btn btn-light border px-4 rounded-3"
            :disabled="loading"
            @click="emit('close')"
          >
            Cancelar
          </button>
          <button
            type="submit"
            class="btn btn-primary px-4 rounded-3 fw-medium"
            :disabled="loading"
          >
            <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status"></span>
            {{ loading ? 'Salvando...' : 'Salvar usuário' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>