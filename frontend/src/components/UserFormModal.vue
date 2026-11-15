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
        if(user){
            name.value = user.name
            email.value = user.email
            login.value = user.login
            active.value = user.active
            currentPassword.value = ''
            newPassword.value = ''
        }  else{
            name.value = ''
            email.value = ''
            login.value = ''
            password.value = ''
            active.value = true
        }
    },
    { immediate: true }
)
async function handleSubmit(){
    loading.value = true
    errorMessage.value = ''

    try{
        if (props.userToEdit){
            const payload = {
                name: name.value,
                email: email.value,
                login: login.value,
                active: active.value
            }
            if (newPassword.value){
            payload.current_password = currentPassword.value
            payload.new_password = newPassword.value
            }
            await userService.update(props.userToEdit.id, payload)
            emit('saved', 'Usuário atualizado com sucesso!')
        } else{
            await userService.create({
                name: name.value,
                email: email.value,
                login: login.value,
                password: password.value,
                active: active.value
            })
            emit('saved', 'Usuário criado com sucesso!')
        }
    } catch(err){
        errorMessage.value = err.message
    } finally{
        loading.value = false
    }
}
</script>


<template>
  <div class="modal-backdrop">
    <div class="modal-box">
      <h2>{{ userToEdit ? 'Editar Usuário' : 'Novo Usuário' }}</h2>

      <p v-if="errorMessage" class="error">{{ errorMessage }}</p>

      <form @submit.prevent="handleSubmit">
        <div>
          <label>Nome:</label>
          <input v-model="name" type="text" required minlength="3" />
        </div>

        <div>
          <label>E-mail:</label>
          <input v-model="email" type="email" required />
        </div>

        <div>
          <label>Login:</label>
          <input v-model="login" type="text" required minlength="3" />
        </div>
        <div v-if="!userToEdit">
          <label>Senha:</label>
          <input v-model="password" type="password" required minlength="8" />
        </div>

        <template v-else>
          <div>
            <label>Senha Atual (para alterar senha):</label>
            <input v-model="currentPassword" type="password" />
          </div>
          <div>
            <label>Nova Senha (mín. 8 caracteres):</label>
            <input v-model="newPassword" type="password" minlength="8" />
          </div>
        </template>

        <div>
          <label>
            <input v-model="active" type="checkbox" />
            Usuário Ativo
          </label>
        </div>

        <div>
          <button type="submit" :disabled="loading">
            {{ loading ? 'Salvando...' : 'Salvar' }}
          </button>
          <button type="button" :disabled="loading" @click="emit('close')">
            Cancelar
          </button>
        </div>
      </form>
    </div>
  </div>
</template>