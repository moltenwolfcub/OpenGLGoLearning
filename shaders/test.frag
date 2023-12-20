#version 330 core
out vec4 FragColor;

in vec2 TexCoord;

void main() {
	FragColor = vec4(TexCoord.x,0.25f,TexCoord.y,1.0f);
}
