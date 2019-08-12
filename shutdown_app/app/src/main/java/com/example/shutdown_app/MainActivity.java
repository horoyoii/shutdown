package com.example.shutdown_app;

import android.support.annotation.NonNull;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;

import android.util.Log;
import android.view.View;
import android.widget.Button;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;
import retrofit2.Retrofit;
import retrofit2.converter.gson.GsonConverterFactory;


// retrofit : https://re-build.tistory.com/14
public class MainActivity extends AppCompatActivity {
    RetrofitExService retrofitExService;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        Retrofit retrofit = new Retrofit.Builder()
                .baseUrl(RetrofitExService.URL)
                .addConverterFactory(GsonConverterFactory.create())
                .build();

        retrofitExService = retrofit.create(RetrofitExService.class);



        Button btn = findViewById(R.id.button);
        btn.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Log.d("HHH", "Clicked");
                //getStatus
                retrofitExService.Shutdown().enqueue(new Callback<Datas>() {
                    @Override
                    public void onResponse(@NonNull Call<Datas> call, @NonNull Response<Datas> response) {
                        if (response.isSuccessful()) {
                            Datas body = response.body();
                            Log.d("HHHH", body.getMessage());

                        }
                    }

                    @Override
                    public void onFailure(@NonNull Call<Datas> call, @NonNull Throwable t) {
                        Log.d("HHH", "Failed");
                        t.printStackTrace();
                    }
                });
            }
        });



    }

}
